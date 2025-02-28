package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/models"
	"google.golang.org/protobuf/encoding/protojson"
)

var cliFramesListCmd = &cobra.Command{
	Use:   "frames-list",
	Short: "List all frames",
	Long: `List all frames: ID, Name,Webhook URL, Domain
This is a wrapper command that uses the farma API.`,
	Run: cliFramesList,
}

func init() {
	rootCmd.AddCommand(cliFramesListCmd)
	cliFramesListCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliFramesListCmd.Flags().String("id", "", "Frame Id or none to list all frames")
	cliFramesListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
func cliFramesList(cmd *cobra.Command, args []string) {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	apiEndpointPath := "frames/"
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

		var res api.ApiResult
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Printf("Failed to parse response: %v", err)
			return
		}
		for _, v := range res.Result {
			item := models.Frame{}
			if err := protojson.Unmarshal(v, &item); err != nil {
				fmt.Printf("Failed to parse user log: %v\n", err)
				continue
			}
			if jsonFlag {
				fmt.Println(string(v))
			} else {
				fmt.Printf("%04d %-32s %45s %s\n", item.Id, item.Name, item.Webhook, item.Domain)
			}
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
