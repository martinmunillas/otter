package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/a-h/templ/cmd/templ/generatecmd"
	"github.com/fsnotify/fsnotify"
	"github.com/martinmunillas/otter/env"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a dev server",
	Long:  `Run a dev server that will auto-compile templ files into go and restart the go server. It will also create a proxy to auto-reload the browser on changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := env.RequiredIntEnvVar("PORT")
		actualPort := port - 1

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

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			runTemplProxy(ctx, port, actualPort)
		}()
		go func() {
			defer wg.Done()
			runReloadServer(ctx, actualPort)
		}()

		wg.Wait()
	},
}

func runTemplProxy(ctx context.Context, port int64, actualPort int64) {
	err := generatecmd.Run(ctx, slog.Default(), generatecmd.Arguments{
		Watch:       true,
		ProxyPort:   int(port),
		Proxy:       fmt.Sprintf("http://localhost:%d", actualPort),
		OpenBrowser: true,
	})
	if err != nil {
		log.Printf("Error running templ command: %v", err)
	}
}

func makeMainCmd(port int64) *exec.Cmd {
	cmd := createDefaultCommand("go", "run", "./cmd/main.go")
	setpgid(cmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	return cmd
}

func runReloadServer(ctx context.Context, port int64) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = addAllGoDirectories(watcher)
	if err != nil {
		log.Fatal(err)
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
			if !event.Has(fsnotify.Write) || !event.Has(fsnotify.Create) || !strings.HasSuffix(event.Name, ".go") {
				continue
			}
			if cmd.Process != nil {
				stop(cmd)
				_ = cmd.Wait()
			}
			log.Printf("Restarting server, file %s changed\n", event.Name)
			cmd = makeMainCmd(port)
			cmd.Start()
			defer stop(cmd)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}

func addAllGoDirectories(w *fsnotify.Watcher) error {
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
			if err != nil {
				return err
			}
		}
		return nil
	})
}
