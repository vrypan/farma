package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliFramesListCmd = &cobra.Command{
	Use:   "frames-list [frameId]",
	Short: "List all frames",
	Long: `List all frames: ID, Name,Webhook URL, Domain
This is a wrapper command that uses the farma API.`,
	Run: cliFramesList,
}

func init() {
	rootCmd.AddCommand(cliFramesListCmd)
	cliFramesListCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliFramesListCmd.Flags().String("start", "", "Start key")
	cliFramesListCmd.Flags().Int("limit", 1000, "Max results")
	cliFramesListCmd.Flags().String("key", "config", "Private key to use")
}
func cliFramesList(cmd *cobra.Command, args []string) {
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	key, _ := cmd.Flags().GetString("key")

	frameId := ""
	if len(args) > 0 {
		frameId = args[0]
	}

	path, _ := cmd.Flags().GetString("path")
	path += "frame/" + frameId

	a := api.ApiClient{}.Init("GET", path, nil, key, frameId)

	next := start
	count := 0
	for {
		if count >= limit {
			break
		}
		resBytes, err := a.Request(next, fmt.Sprintf("%d", limit))
		if err != nil {
			fmt.Println(string(resBytes))
			fmt.Printf("Failed to make API call: %v\n", err)
			return
		}
		var res api.ApiResult
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Println(string(resBytes))
			fmt.Printf("Failed to parse response: %v\n", err)
			return
		}
		for _, v := range res.Result {
			fmt.Println(string(v))
		}
		if res.Next == "" {
			break
		}
		next = res.Next
		count += len(res.Result)
	}
}
