package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/martinmunillas/otter/log"
	"github.com/martinmunillas/otter/utils"
	"github.com/spf13/cobra"
)

func refactorFunc(old, new string, filePatterns []string, logger *slog.Logger) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		var matched bool
		for _, pattern := range filePatterns {
			var err error
			matched, err = filepath.Match(pattern, fi.Name())
			if err != nil {
				return err
			}
			if matched {
				read, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				logger.Debug(fmt.Sprintf("Refactoring: %s", path))

				newContents := strings.Replace(string(read), old, new, -1)

				err = os.WriteFile(path, []byte(newContents), 0)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Otter project",
	Long:  `Set up a new Otter project to get you developing asap`,
	Run: func(c *cobra.Command, args []string) {
		if len(args) < 1 {
			panic("You should provide a project name in the form `otter init {githubUser}/{repositoryName}`")
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := log.NewLogger(verbose)
		projectName := args[0]
		githubUser := ""
		repositoryName := ""
		collectingRepository := false
		for _, c := range projectName {
			if c == '/' {
				collectingRepository = true
				continue
			}
			if collectingRepository {
				repositoryName += string(c)
			} else {
				githubUser += string(c)
			}
		}
		if githubUser == "" || repositoryName == "" {
			fatal(logger, errors.New("you should provide a project name in the form `otter init {githubUser}/{repositoryName}`"))
		}
		err := createDefaultCommand("git", "clone", "https://github.com/martinmunillas/otter-example", repositoryName).Run()
		if err != nil {
			fatal(logger, fmt.Errorf("error cloning example project: %v", err))
		}

		// ignore the error because if the repositoryName is invalid git clone would have failed already
		cwd, _ := filepath.Abs(fmt.Sprintf("./%s", repositoryName))

		err = os.RemoveAll(fmt.Sprintf("%s/.git/", repositoryName))
		if err != nil {
			fatal(logger, fmt.Errorf("error recreating git repository: %v", err))
		}
		cmd := createDefaultCommand("git", "init")
		cmd.Dir = cwd
		err = cmd.Run()
		if err != nil {
			fatal(logger, fmt.Errorf("error recreating git repository: %v", err))
		}
		err = filepath.Walk(cwd, refactorFunc("martinmunillas/otter-example", projectName, []string{"*.go", "*.templ", "go.mod", "go.sum"}, logger))
		if err != nil {
			fatal(logger, fmt.Errorf("error refactoring go module: %v", err))
		}
		cmd = createDefaultCommand("go", "generate")
		cmd.Dir = cwd
		err = cmd.Run()
		if err != nil {
			fatal(logger, fmt.Errorf("error generating templ files: %v", err))
		}
		err = utils.CopyFile(fmt.Sprintf("%s/.env.example", cwd), fmt.Sprintf("%s/.env", cwd))
		if err != nil {
			fatal(logger, fmt.Errorf("error setting up .env file: %v", err))
		}

	},
}
