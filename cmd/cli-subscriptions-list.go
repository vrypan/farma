package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/models"
	"google.golang.org/protobuf/encoding/protojson"
)

var cliSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "List all logs",
	Long: `List all user logs: UserId (Fid), FrameId, AppId (Fid), Event, Timestamp
This is a wrapper command that uses the farma API.`,
	Run: cliSubscriptions,
}

func init() {
	rootCmd.AddCommand(cliSubscriptionsCmd)
	cliSubscriptionsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliSubscriptionsCmd.Flags().String("id", "", "User Id (fid) or none to list all subscriptions")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	apiEndpointPath := "subscriptions/"
	endpoint, _ := cmd.Flags().GetString("path")
	id, _ := cmd.Flags().GetString("id")
	body := []byte("")
	method := "GET"
	next := ""
	for {
		resBytes, err := api.ApiCall(method, endpoint, apiEndpointPath, id, body, "start="+next)
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
			item := models.Subscription{}
			if err := protojson.Unmarshal(v, &item); err != nil {
				fmt.Printf("Failed to parse user log: %v\n", err)
				continue
			}
			fmt.Printf("%06d %04d %06d %-20s %s %s %s %s\n",
				item.UserId, item.FrameId, item.AppId, item.Status.String(),
				item.Ctime.AsTime().Format(time.RFC3339), item.Mtime.AsTime().Format(time.RFC3339), item.Token, item.Url,
			)
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
