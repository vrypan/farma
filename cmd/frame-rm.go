package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var frameRmCmd = &cobra.Command{
	Use:     "remove [ID]",
	Short:   "Remove a frame",
	Aliases: []string{"rm"},
	Run:     rmFrame,
}

func init() {
	frameCmd.AddCommand(frameRmCmd)
}
func rmFrame(cmd *cobra.Command, args []string) {
	fmt.Println("Frame deletion not implemented yet.")
}
