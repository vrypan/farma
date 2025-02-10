package cmd

import (
	"github.com/spf13/cobra"
)

var frameCmd = &cobra.Command{
	Use:   "frame",
	Short: "Add/Remove/List frames",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(frameCmd)
}
