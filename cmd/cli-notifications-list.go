package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliNotificationsListCmd = &cobra.Command{
	Use:   "notifications-list [frameId] [notificationId]",
	Short: "List notifications",
	Run:   cliNotifications,
}

func init() {
	rootCmd.AddCommand(cliNotificationsListCmd)
	cliNotificationsListCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliNotificationsListCmd.Flags().String("start", "", "Start key")
	cliNotificationsListCmd.Flags().Int("limit", 1000, "Max results")
	cliNotificationsListCmd.Flags().String("key", "config", "Private key to use")

}

func cliNotifications(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")

	path, _ := cmd.Flags().GetString("path")
	path += "notification/"
	if len(args) < 1 {
		cmd.Help()
	}
	path += args[0]
	if len(args) > 1 {
		path += "/" + args[1]
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
