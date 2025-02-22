package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
)

var cliNotificationSendCmd = &cobra.Command{
	Use:   "notification-send [frame] [title] [body] [url]",
	Short: "Configure a new frame",
	Run:   cliNotificationSend,
}

func init() {
	rootCmd.AddCommand(cliNotificationSendCmd)
	cliNotificationSendCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
}

func cliNotificationSend(cmd *cobra.Command, args []string) {
	apiEndpointPath := "notifications/"
	endpoint, _ := cmd.Flags().GetString("path")
	method := "POST"

	if len(args) != 4 {
		cmd.Help()
		return
	}
	frameId := args[0]
	title := args[1]
	body := args[2]
	url := args[3]

	payload := `{
		"frameId": ` + frameId + `,
		"title": "` + title + `",
		"body": "` + body + `",
		"url": "` + url + `"
	}`

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, "", []byte(payload))
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
