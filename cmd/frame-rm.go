package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
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
	db.Open()
	defer db.Close()

	if len(args) == 0 {
		fmt.Fprintln(cmd.OutOrStderr(), "Error: Frame ID not defined")
		os.Exit(1)
	}
	id := args[0]

	sql := fmt.Sprintf("DELETE FROM frames WHERE id=%s", id)
	_, err := db.Instance.Exec(sql)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error:", err)
		os.Exit(1)
	}
	fmt.Println("Frame deleted")
}
