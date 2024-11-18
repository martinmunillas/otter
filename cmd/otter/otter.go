package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "otter",
	Short: "Otter is a toolkit library for building go and templ applications",
	Long:  ``,
}

func main() {
	_ = godotenv.Load()
	cmd.AddCommand(devCmd)
	cmd.AddCommand(initCmd)
	cmd.AddCommand(migrateCmd)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
