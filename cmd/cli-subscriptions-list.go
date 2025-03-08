package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions-list [frameId]",
	Short: "List subscriptions",
	Run:   cliSubscriptions,
}

func init() {
	rootCmd.AddCommand(cliSubscriptionsCmd)
	cliSubscriptionsCmd.Flags().String("path", "", "API endpoint.")
	cliSubscriptionsCmd.Flags().String("start", "", "Start key")
	cliSubscriptionsCmd.Flags().Int("limit", 1000, "Max results")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	path := "/api/v2/subscription/"
	if len(args) > 0 {
		path += args[0]
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
