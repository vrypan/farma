package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
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
	cliSubscriptionsCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliSubscriptionsCmd.Flags().String("id", "", "User Id (fid) or none to list all subscriptions")
}

func cliSubscriptions(cmd *cobra.Command, args []string) {
	apiEndpointPath := "subscriptions/"
	endpoint, _ := cmd.Flags().GetString("path")
	id, _ := cmd.Flags().GetString("id")
	body := []byte("")
	method := "GET"

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, id, body)
	if err != nil {
		fmt.Printf("Failed to make API call: %v", err)
		return
	}

	fmt.Println(string(res))
	dataStruct := struct {
		Data []*utils.Subscription `json:"subscriptions"`
	}{}

	if err := json.Unmarshal(res, &dataStruct); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	for _, data := range dataStruct.Data {
		fmt.Printf("%06d %04d %06d %-20s %s %s %s %s\n",
			data.UserId, data.FrameId, data.AppId, data.Status.String(),
			data.Ctime.AsTime().Format(time.RFC3339), data.Mtime.AsTime().Format(time.RFC3339), data.Token, data.Url,
		)
	}
}
