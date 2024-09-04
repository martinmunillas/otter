package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/a-h/templ/cmd/templ/generatecmd"
	"github.com/fsnotify/fsnotify"
	"github.com/martinmunillas/otter/env"
	"github.com/martinmunillas/otter/log"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a dev server",
	Long:  `Run a dev server that will auto-compile templ files into go and restart the go server. It will also create a proxy to auto-reload the browser on changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := env.RequiredIntEnvVar("PORT")
		actualPort := port - 1

		verbose, _ := cmd.Flags().GetBool("verbose")

		userInterrupt := make(chan os.Signal, 1)
		signal.Notify(userInterrupt, syscall.SIGTERM, syscall.SIGINT)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			<-userInterrupt
			cancel()
			// TODO: find out how to get rid of this sleep, right now is necessary because we need the cancel signal to be sent to the chanel so they can stop their subprocesses and only exit after that, otherwise orphan processes stay around.
			time.Sleep(time.Second / 2)
			os.Exit(1)
		}()

		logger := log.NewLogger(verbose)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			runTemplProxy(ctx, logger, port, actualPort)
		}()
		go func() {
			defer wg.Done()
			runReloadServer(ctx, logger, actualPort)
		}()

		wg.Wait()
	},
}

func init() {
	devCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

func runTemplProxy(ctx context.Context, logger *slog.Logger, port int64, actualPort int64) {
	err := generatecmd.Run(ctx, logger, generatecmd.Arguments{
		Watch:             true,
		ProxyPort:         int(port),
		Proxy:             fmt.Sprintf("http://localhost:%d", actualPort),
		OpenBrowser:       true,
		KeepOrphanedFiles: false,
		WorkerCount:       runtime.NumCPU(),
		IncludeVersion:    true,
		Path:              ".",
	})
	if err != nil {
		logger.Error(err.Error())
	}
}

func makeMainCmd(port int64) *exec.Cmd {
	cmd := createDefaultCommand("go", "run", "./cmd/main.go")
	setpgid(cmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	return cmd
}

func runReloadServer(ctx context.Context, logger *slog.Logger, port int64) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error(err.Error())
	}
	defer watcher.Close()

	err = addAllGoDirectories(watcher, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	cmd := makeMainCmd(port)
	cmd.Start()

	for {
		select {
		case <-ctx.Done():
			stop(cmd)
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !(event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) || !strings.HasSuffix(event.Name, ".go") {
				continue
			}
			logger.Debug(fmt.Sprintf("Go file %s changed", event.Name))
			if cmd.Process != nil {
				logger.Debug("Stopping running server")
				stop(cmd)
				_ = cmd.Wait()
				logger.Debug("Stopped running server")
			}
			logger.Debug("Restarting server")
			cmd = makeMainCmd(port)
			cmd.Start()
			defer stop(cmd)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error(err.Error())
		}
	}

}

func addAllGoDirectories(w *fsnotify.Watcher, logger *slog.Logger) error {
	root, err := filepath.Abs("./")
	if err != nil {
		return err
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			dir := filepath.Dir(path)
			err = w.Add(dir)
			logger.Debug(fmt.Sprintf("Watching %s directory", dir))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
