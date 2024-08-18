package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/martinmunillas/otter/env"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
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
}
