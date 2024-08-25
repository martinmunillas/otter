package main

import (
	"fmt"
	"io/fs"
	"log"
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

		var wg sync.WaitGroup
		wg.Add(2)

		go runTemplProxy(&wg, port)
		go runReloadServer(&wg)

		wg.Wait()
	},
}

func runTemplProxy(wg *sync.WaitGroup, port int64) {
	defer wg.Done()
	actualPort := port - 1
	err := createDefaultCommand("templ", "generate", "--watch", fmt.Sprintf("--proxy=http://localhost:%d", actualPort), fmt.Sprintf("--proxyport=%d", port)).Run()
	if err != nil {
		log.Printf("Error running templ command: %v", err)
	}

}

func runReloadServer(wg *sync.WaitGroup) {
	defer wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	addAllGoDirectories(watcher)

	cmd := createDefaultCommand("go", "run", "./cmd/main.go")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
				cmd = createDefaultCommand("go", "run", "./cmd/main.go")
				err = cmd.Run()
				if err != nil {
					log.Fatalf(err.Error())
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}

}

func addAllGoDirectories(w *fsnotify.Watcher) error {
	return filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			dir := filepath.Dir(path)
			err := w.Add(dir)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
