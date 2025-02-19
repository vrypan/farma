package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/utils"
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
	cliSubscriptionsCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	payload := `{
		"command": "subscriptions/get",
		"params": {}
	}`

	cmd.Flags().Bool("send", true, "")
	_, resp := SendCommand(cmd, payload)
	if resp == nil {
		return
	}

	var j map[string]interface{}
	if err := json.Unmarshal(resp, &j); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}

	switch status := j["status"].(string); status {
	case "SUCCESS":
		list, ok := j["data"].([]interface{})
		if !ok {
			fmt.Println("Failed to cast data to list of map[string]interface{}")
			return
		}
		for _, data := range list {
			d, ok := data.(map[string]interface{})
			if !ok {
				fmt.Println("Failed to cast list item to map[string]interface{}")
				continue
			}
			userId, _ := d["userId"].(float64)
			frameId, _ := d["frameId"].(float64)
			appId, _ := d["appId"].(float64)
			statusNum, _ := d["status"].(float64)
			url, _ := d["url"].(string)
			token, _ := d["token"].(string)
			status := utils.SubscriptionStatus(statusNum).String()
			cTime, _ := d["ctime"].(map[string]interface{})
			mTime, _ := d["mtime"].(map[string]interface{})
			cdate := time.Unix(int64(cTime["seconds"].(float64)), int64(cTime["nanos"].(float64)))
			mdate := time.Unix(int64(mTime["seconds"].(float64)), int64(mTime["nanos"].(float64)))

			fmt.Printf("%04d %06d %06d %-20s %s %s %s %s\n",
				int(frameId), int(userId), int(appId), status, cdate.Format("2006-01-02 15:04:05"), mdate.Format("2006-01-02 15:04:05"), token, url)
		}
	case "ERROR":
		message, _ := j["message"].(string)
		fmt.Printf("An error occurred: %s\n", message)
	default:
		fmt.Printf("Unknown status received: %s\n", status)
	}
}
