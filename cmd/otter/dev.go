package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

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

		var wg sync.WaitGroup
		wg.Add(2)

		go runTemplProxy(&wg, port, actualPort)
		go runReloadServer(&wg, actualPort)

		wg.Wait()
	},
}

func runTemplProxy(wg *sync.WaitGroup, port int64, actualPort int64) {
	defer wg.Done()
	cmd := createDefaultCommand("templ", "generate", "--watch", fmt.Sprintf("--proxy=http://localhost:%d", actualPort), fmt.Sprintf("--proxyport=%d", port))
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running templ command: %v", err)
	}
	defer cmd.Process.Kill()
}

func makeMainCmd(port int64) *exec.Cmd {
	cmd := createDefaultCommand("go", "run", "./cmd/main.go")
	setpgid(cmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	return cmd
}

func runReloadServer(wg *sync.WaitGroup, port int64) {
	defer wg.Done()

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
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) && strings.HasSuffix(event.Name, ".go") {
				println(cmd.Process)
				if cmd.Process != nil {
					stop(cmd)
					_ = cmd.Wait()
				}
				log.Printf("Restarting server, file %s changed\n", event.Name)
				cmd = makeMainCmd(port)
				cmd.Start()
				defer stop(cmd)
			}
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
