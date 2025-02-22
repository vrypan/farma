package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/utils"
)

var cliFramesListCmd = &cobra.Command{
	Use:   "frames-list",
	Short: "List all frames",
	Long: `List all frames: ID, Name,Webhook URL, Domain
This is a wrapper command that uses the farma API.`,
	Run: cliFramesList,
}

func init() {
	rootCmd.AddCommand(cliFramesListCmd)
	cliFramesListCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliFramesListCmd.Flags().String("id", "", "Frame Id or none to list all frames")
}
func cliFramesList(cmd *cobra.Command, args []string) {
	apiEndpointPath := "frames/"
	endpoint, _ := cmd.Flags().GetString("path")
	id, _ := cmd.Flags().GetString("id")
	body := []byte("")
	method := "GET"

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, id, body)
	if err != nil {
		fmt.Printf("Failed to make API call: %v", err)
		return
	}

	var list []*utils.Frame
	if err := json.Unmarshal(res, &list); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	for _, item := range list {
		fmt.Printf("%04d %-32s %45s %s\n", item.Id, item.Name, item.Webhook, item.Domain)
	}

}
