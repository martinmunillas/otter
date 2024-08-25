package main

import (
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
	cmd.AddCommand(devCmd)
	cmd.AddCommand(initCmd)
	cmd.Execute()

}
