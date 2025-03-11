package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliNotificationSendCmd = &cobra.Command{
	Use:   "notification-send [frameId] [title] [body] [url] [userId1] [userId2] ...",
	Short: "Send notifications",
	Run:   cliNotificationSend,
}

func init() {
	rootCmd.AddCommand(cliNotificationSendCmd)
	cliNotificationSendCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v2/frames/ (from config file)")
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
