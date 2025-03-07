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
	Use:   "user-logs [frameId]",
	Short: "List all logs",
	Long: `List user logs: UserId (Fid), FrameId, AppId (Fid), Event, Timestamp
This is a wrapper command that uses the farma API.`,
	Run: cliLogs,
}

func init() {
	rootCmd.AddCommand(cliLogsCmd)
	cliLogsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliLogsCmd.Flags().String("id", "", "User Id (fid) or none to list all logs")
	cliLogsCmd.Flags().BoolP("json", "j", false, "Output logs in JSON format")
	cliLogsCmd.Flags().Int("frame", 0, "Frame id")
}

func cliLogs(cmd *cobra.Command, args []string) {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	userId, _ := cmd.Flags().GetString("id")

	var id string
	if len(args) > 0 {
		id = args[0]
	}

	a := api.ApiCallData{}
	a.Path = "logs/" + id
	if len(args) != 0 {
		a.Path += args[0]
	}
	if userId != "" {
		a.Path += "/" + userId
	}
	a.Endpoint, _ = cmd.Flags().GetString("path")
	a.Body = ""
	a.Method = "GET"

	next := ""
	for {
		if next != "" {
			a.RawQuery = fmt.Sprintf("start=%s", next)
		}
		resBytes, err := a.Call()
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
			if jsonFlag {
				fmt.Println(string(v))
			} else {
				fmt.Printf("%06d %04d %04d %-45s %s",
					item.UserId, item.FrameId, item.AppId, item.EvtType.Enum(), item.Ctime.AsTime().Format(time.RFC3339))
				if item.GetEventContextNotification() != nil {
					fmt.Printf(" %s", item.GetEventContextNotification().GetId())
				}
				fmt.Println()
			}
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
