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
	Use:   "subscriptions-list [frameId]",
	Short: "List subscriptions",
	Long: `List frame subscriptions.
Fields: UserId (Fid), FrameId, AppId (Fid), Status, Ctime, Mtime, Token, Endpoint
If frameId is provided, show only this frame.
This is a wrapper command that uses the farma API.`,
	Run: cliSubscriptions,
}

func init() {
	rootCmd.AddCommand(cliSubscriptionsCmd)
	cliSubscriptionsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliSubscriptionsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	var id string
	if len(args) > 0 {
		id = args[0]
	}

	a := api.ApiCallData{}
	a.Path = "subscriptions/" + id
	if len(args) != 0 {
		a.Path += args[0]
	}
	a.Endpoint, _ = cmd.Flags().GetString("path")
	a.Body = ""
	a.Method = "GET"

	next := ""
	for {
		if next != "" {
			a.RawQuery = fmt.Sprintf("start=%s", next)
		}
		resBytes, err := a.Call()
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
			if jsonFlag {
				fmt.Println(string(v))
			} else {
				fmt.Printf("%06d %04d %06d %-20s %s %s %s %s\n",
					item.UserId, item.FrameId, item.AppId, item.Status.String(),
					item.Ctime.AsTime().Format(time.RFC3339), item.Mtime.AsTime().Format(time.RFC3339), item.Token, item.Url,
				)
			}
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
