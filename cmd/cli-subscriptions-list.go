package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions-list frame_id",
	Short: "List subscriptions",
	Long: `List all subscriptions for frame_id. It will show active and
inaactive subscriptions (users that disabled notifications or
removed the frame)`,
	Run: cliSubscriptions,
}

func init() {
	rootCmd.AddCommand(cliSubscriptionsCmd)
	cliSubscriptionsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v2/ (host.addr from config)")
	cliSubscriptionsCmd.Flags().String("start", "", "Start key")
	cliSubscriptionsCmd.Flags().Int("limit", 1000, "Max results")
	cliSubscriptionsCmd.Flags().String("key", "config", "Private key to use")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")

	path, _ := cmd.Flags().GetString("path")
	path += "subscription/"
	if len(args) > 0 {
		path += args[0]
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
