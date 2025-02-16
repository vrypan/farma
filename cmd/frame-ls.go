package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var frameLsCmd = &cobra.Command{
	Use:     "list [frame name",
	Short:   "List frames",
	Aliases: []string{"ls"},
	Run:     lsFrame,
}

func init() {
	frameCmd.AddCommand(frameLsCmd)
	frameLsCmd.Flags().BoolP("with-subscriptions", "", false, "Show subscriptions for each frame")
}

func lsFrame(cmd *cobra.Command, args []string) {
	//subscriptionsFlag, _ := cmd.Flags().GetBool("with-subscriptions")

	db.Open()
	defer db.Close()

	frames := utils.AllFrames()
	for _, frame := range frames {
		fmt.Printf("%04d %-32s %s %s\n", frame.Id, frame.Name, frame.Endpoint, frame.Domain)
	}
}
