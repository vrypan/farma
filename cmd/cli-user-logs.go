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

	next := ""
	for {
		resBytes, err := api.ApiCall(method, endpoint, apiEndpointPath, id, body, "start="+next)
		if err != nil {
			fmt.Printf("Failed to make API call: %v", err)
			return
		}
		type apiResult struct {
			Error  string            `json:"error"`
			Result []json.RawMessage `json:"result"`
			Next   string            `json:"next"`
		}
		var res apiResult
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Printf("Failed to parse response: %v\n", err)
			return
		}

		for _, v := range res.Result {
			item := models.UserLog{}
			if err := protojson.Unmarshal(v, &item); err != nil {
				fmt.Printf("Failed to parse user log: %v\n", err)
				continue
			}
			fmt.Printf("%06d %04d %04d %-45s %s",
				item.UserId, item.FrameId, item.AppId, item.EvtType.Enum(), item.Ctime.AsTime().Format(time.RFC3339))
			if item.GetEventContextNotification() != nil {
				fmt.Printf(" %s", item.GetEventContextNotification().GetId())
			}
			fmt.Println()
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
