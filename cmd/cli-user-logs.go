package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliLogsCmd = &cobra.Command{
	Use:   "user-logs frame_id [fid]",
	Short: "User logs",
	Long: `Detailed user logs for a frame: frame additions/removals,
notifications enabled/disabled, notifications sent to each user.`,
	Run: cliLogs,
}

func init() {
	rootCmd.AddCommand(cliLogsCmd)
	cliLogsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v2/ (host.addr from config)")
	cliLogsCmd.Flags().String("start", "", "Start key")
	cliLogsCmd.Flags().Int("limit", 1000, "Max results")
	cliLogsCmd.Flags().String("key", "config", "Private key to use")

}

func cliLogs(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	path, _ := cmd.Flags().GetString("path")
	path += "logs/"
	if len(args) > 0 {
		path += args[0] + "/"
	}
	if len(args) > 1 {
		path += args[1]
	}

	a := api.ApiClient{}.Init("GET", path, nil, key, "")
	next := start
	count := 0
	for {
		if count >= limit {
			break
		}
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
		count += len(res.Result)
	}
}
