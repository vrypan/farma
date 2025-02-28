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
	cliNotificationsListCmd.Flags().BoolP("json", "j", false, "Display entries as JSON")

}

func cliNotifications(cmd *cobra.Command, args []string) {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	apiEndpointPath := "notifications/"
	endpoint, _ := cmd.Flags().GetString("path")
	notificationId, _ := cmd.Flags().GetInt("notificationid")
	var id string
	if notificationId == 0 {
		id = ""
	} else {
		id = fmt.Sprintf("%d", notificationId)
	}
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
			item := models.Notification{}
			if err := protojson.Unmarshal(v, &item); err != nil {
				fmt.Printf("Failed to parse user log: %v\n", err)
				continue
			}

			if jsonFlag == true {
				fmt.Println(string(v))
			} else {
				fmt.Printf("%s %s %3d %3d %3d %s %s %s\n",
					item.Id, item.Ctime.AsTime().Format(time.RFC3339),
					len(item.GetSuccessTokens()), len(item.GetFailedTokens()), len(item.GetRateLimitedTokens()),
					item.Title, item.Message, item.Link,
				)
			}

		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
