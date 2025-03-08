package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliNotificationSendCmd = &cobra.Command{
	Use:   "notification-send [frameId] [title] [body] [url] [userId1] [userId2] ...",
	Short: "Send notifications",
	Long: `Send a notification to <frameId> subscribers.
If <url> is "", the notification link will open the frame.
If a list of userIds is provided, the notification will
be sent to all of them, provided they are subscribed to <formId>`,
	Run: cliNotificationSend,
}

func init() {
	rootCmd.AddCommand(cliNotificationSendCmd)
	cliNotificationSendCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
}

func cliNotificationSend(cmd *cobra.Command, args []string) {
	if len(args) < 4 {
		cmd.Help()
		return
	}

	userIdsStr := ""
	for i := 4; i <= len(args)-1; i++ {
		userIdsStr += args[i]
		if i < len(args)-1 {
			userIdsStr += ","
		}
	}
	payload := `{
		"frameId": ` + args[0] + `,
		"title": "` + args[1] + `",
		"body": "` + args[2] + `",
		"url": "` + args[3] + `",
		"userIds": [` + userIdsStr + `]
	}`
	a := api.ApiClient{}.Init("POST", "notification/", []byte(payload), []byte("config"), "")
	res, err := a.Request("", "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	fmt.Println(string(res))
}
