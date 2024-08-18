package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := runCommand("templ", "generate", "--watch", fmt.Sprintf("--proxy=http://localhost:%s", os.Getenv("PORT")), "-v")
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
