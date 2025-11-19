package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = cobra.Command{
	Use:   "template",
	Short: "All Commands",
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
