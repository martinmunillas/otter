package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "otter",
	Short: "Otter is a toolkit library for building go and templ applications",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cmd.AddCommand(devCmd)
	cmd.AddCommand(initCmd)
	cmd.AddCommand(migrateCmd)
	cmd.Execute()

}
