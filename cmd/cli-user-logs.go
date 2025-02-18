package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/utils"
)

var cliLogsCmd = &cobra.Command{
	Use:   "user-logs",
	Short: "List all logs",
	Long: `List all user logs: UserId (Fid), FrameId, AppId (Fid), Event, Timestamp
This is a wrapper command that uses the farma API.`,
	Run: cliLogs,
}

func init() {
	rootCmd.AddCommand(cliLogsCmd)
	cliLogsCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
}

func cliLogs(cmd *cobra.Command, args []string) {
	payload := `{
		"command": "logs/get",
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
			evtTypeNum, _ := d["evtType"].(float64)
			evtType := utils.EventType(evtTypeNum).String()
			cTime, _ := d["ctime"].(map[string]interface{})
			date := time.Unix(int64(cTime["seconds"].(float64)), int64(cTime["nanos"].(float64)))

			fmt.Printf("%06d %04d %04d %-45s %s\n", int(userId), int(frameId), int(appId), evtType, date.Format("2006-01-02 15:04:05"))
		}
	case "ERROR":
		message, _ := j["message"].(string)
		fmt.Printf("An error occurred: %s\n", message)
	default:
		fmt.Printf("Unknown status received: %s\n", status)
	}
}
