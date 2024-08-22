package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/martinmunillas/otter/env"
	"github.com/spf13/cobra"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

var cmd = &cobra.Command{
	Use:   "otter",
	Short: "Otter is a toolkit library for building go and templ applications",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a dev server",
	Long:  `Run a dev server that will auto-compile templ files into go and restart the go server. It will also create a proxy to auto-reload the browser on changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := env.RequiredStringEnvVar("PORT")

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			err := runCommand("templ", "generate", "--watch", fmt.Sprintf("--proxy=http://localhost:%s", port), "-v")
			if err != nil {
				log.Printf("Error running templ command: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			err := runCommand("wgo", "run", "./cmd")
			if err != nil {
				log.Printf("Error running server: %v", err)
			}
		}()

		wg.Wait()
	},
}

func main() {
	cmd.AddCommand(devCmd)
	cmd.Execute()

}
