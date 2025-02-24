package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/models"
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
	cliLogsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliLogsCmd.Flags().String("id", "", "User Id (fid) or none to list all logs")
}

func cliLogs(cmd *cobra.Command, args []string) {
	apiEndpointPath := "logs/"
	endpoint, _ := cmd.Flags().GetString("path")
	id, _ := cmd.Flags().GetString("id")
	body := []byte("")
	method := "GET"

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, id, body)
	if err != nil {
		fmt.Printf("Failed to make API call: %v", err)
		return
	}

	var list []*models.UserLog
	if err := json.Unmarshal(res, &list); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	for _, item := range list {
		fmt.Printf("%06d %04d %04d %-45s %s\n",
			item.UserId, item.FrameId, item.AppId, item.EvtType.Enum(), item.Ctime.AsTime().Format(time.RFC3339))
	}
}
