package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var cliNotificationSendCmd = &cobra.Command{
	Use:   "notification-send [frame] [title] [body] [url]",
	Short: "Configure a new frame",
	Run:   cliNotificationSend,
}

func init() {
	rootCmd.AddCommand(cliNotificationSendCmd)
	cliNotificationSendCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
}

func cliNotificationSend(cmd *cobra.Command, args []string) {
	cmd.Flags().Bool("send", true, "")

	if len(args) != 4 {
		cmd.Help()
		return
	}
	frameName := args[0]
	title := args[1]
	body := args[2]
	url := args[3]

	payload := `{
		"command": "notification/send",
		"params": {
			"frame": "` + frameName + `",
			"title": "` + title + `",
			"body": "` + body + `",
			"url": "` + url + `"
		}
	}`
	_, resp := SendCommand(cmd, payload)

	if resp != nil {
		var j map[string]any
		json.Unmarshal(resp, &j)
		if j["status"].(string) == "SUCCESS" {
			data := j["data"].(map[string]any)
			fmt.Printf("Notification Id: %s\n", data["NotificationId"].(string))
			fmt.Printf("Token count: %d\n", int(data["Count"].(float64)))
		}
		if j["status"].(string) == "ERROR" {
			fmt.Printf("%s: %s\n", j["message"].(string), j["data"].(string))
		}
	}

}
