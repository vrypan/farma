package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
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
	cliFramesListCmd.Flags().String("start", "", "Start key")
	cliFramesListCmd.Flags().Int("limit", 1000, "Max results")
}
func cliFramesList(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	path, _ := cmd.Flags().GetString("path")

	if path == "" {
		path = "/api/v2/frame/"
	}
	a := api.ApiClient{}.Init("GET", path+id, nil, []byte("config"), "")

	next := start
	for {
		resBytes, err := a.Request(next, fmt.Sprintf("%d", limit))
		if err != nil {
			fmt.Printf("Failed to make API call: %v", err)
			return
		}
		var res api.ApiResult
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Printf("Failed to parse response: %v", err)
			return
		}
		for _, v := range res.Result {
			fmt.Println(string(v))
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
