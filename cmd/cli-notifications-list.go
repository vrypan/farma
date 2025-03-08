package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliNotificationsListCmd = &cobra.Command{
	Use:   "notifications-list [notificationId]",
	Short: "List notifications",
	Long: `List frame notifications.
Fields: UserId (Fid), FrameId, AppId (Fid), Status, Ctime, Mtime, Token, Endpoint
Optionally, show only one notification if notificationId is provided`,
	Run: cliNotifications,
}

func init() {
	rootCmd.AddCommand(cliNotificationsListCmd)
	cliNotificationsListCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliNotificationsListCmd.Flags().String("start", "", "Start key")
	cliNotificationsListCmd.Flags().Int("limit", 1000, "Max results")
}

func cliNotifications(cmd *cobra.Command, args []string) {
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	path := "/api/v2/notification/"
	if len(args) > 0 {
		path = fmt.Sprintf("/api/v2/notification/%s", args[0])
	}
	a := api.ApiClient{}.Init("GET", path, nil, []byte("config"), "")
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
