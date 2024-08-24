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

func runCommand(env []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(cmd.Env, env...)
	return cmd.Run()
}

var cmd = &cobra.Command{
	Use:   "otter",
	Short: "Otter is a toolkit library for building go and templ applications",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs a dev server",
	Long:  `Run a dev server that will auto-compile templ files into go and restart the go server. It will also create a proxy to auto-reload the browser on changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := env.RequiredIntEnvVar("PORT")
		actualPort := port - 1

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			err := runCommand(nil, "templ", "generate", "--watch", fmt.Sprintf("--proxy=http://localhost:%d", actualPort), fmt.Sprintf("--proxyport=%d", port))
			if err != nil {
				log.Printf("Error running templ command: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			err := runCommand([]string{fmt.Sprintf("PORT=%d", actualPort)}, "wgo", "run", "./cmd")
			if err != nil {
				log.Printf("Error running server: %v", err)
			}
		}()

		wg.Wait()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Otter project",
	Long:  `Set up a new Otter project to get you developing asap`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic("You should provide a project name in the form `otter init {project-name}`")
		}
		err := runCommand(nil, "git", "clone", "https://github.com/martinmunillas/otter-example", args[0])
		if err != nil {
			log.Printf("Error cloning example project: %v", err)
		}
	},
}

func main() {
	cmd.AddCommand(devCmd)
	cmd.AddCommand(initCmd)
	cmd.Execute()

}
