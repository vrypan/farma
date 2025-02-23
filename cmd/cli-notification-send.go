package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
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
	apiEndpointPath := "notifications/"
	endpoint, _ := cmd.Flags().GetString("path")
	method := "POST"

	frameId := args[0]
	title := args[1]
	body := args[2]
	url := args[3]

	userIdsStr := ""
	for i := 4; i <= len(args)-1; i++ {
		userIdsStr += args[i]
		if i < len(args)-1 {
			userIdsStr += ","
		}
	}
	payload := `{
		"frameId": ` + frameId + `,
		"title": "` + title + `",
		"body": "` + body + `",
		"url": "` + url + `",
		"userIds": [` + userIdsStr + `]
	}`

	fmt.Println("DEBUG. Payload=", payload)
	res, err := api.ApiCall(method, endpoint, apiEndpointPath, "", []byte(payload))
	fmt.Println("DEBUG, Response=", string(res))
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	var data struct {
		NotificationId string
		Count          int
	}

	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Failed to parse response: %v", err)

		return
	}
	fmt.Printf("NotificationId:%s, Count:%d\n", data.NotificationId, data.Count)

}
