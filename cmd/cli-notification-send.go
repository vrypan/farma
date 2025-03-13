package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliNotificationSendCmd = &cobra.Command{
	Use:   "notification-send frame_id title body url [fid1] [fid2] ...",
	Short: "Send a notification",
	Long: `Send a notification to frame_id subscribers. Title, Body and Url
must be specificed. Set Url to "" for the frame's domain.
If a list of fids is provided, only these users will be notified, provided
they have active subscriptions.`,
	Run: cliNotificationSend,
}

func init() {
	rootCmd.AddCommand(cliNotificationSendCmd)
	cliNotificationSendCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v2/ (host.addr from config)")
	cliNotificationSendCmd.Flags().String("key", "config", "Private key to use")
}

func cliNotificationSend(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	if len(args) < 4 {
		cmd.Help()
		return
	}

	frameId := args[0]
	path, _ := cmd.Flags().GetString("path")
	path += "notification/" + frameId

	userIdsStr := ""
	for i := 4; i <= len(args)-1; i++ {
		userIdsStr += args[i]
		if i < len(args)-1 {
			userIdsStr += ","
		}
	}
	payload := `{
		"frameId": "` + frameId + `",
		"title": "` + args[1] + `",
		"body": "` + args[2] + `",
		"url": "` + args[3] + `",
		"userIds": [` + userIdsStr + `]
	}`

	a := api.ApiClient{}.Init("POST", path, []byte(payload), key, "")
	res, err := a.Request("", "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	fmt.Println(string(res))
}
