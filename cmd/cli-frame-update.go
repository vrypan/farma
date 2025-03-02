package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/models"
)

var cliFrameUpdCmd = &cobra.Command{
	Use:   "frame-update [id] [name] [domain]",
	Short: "Update a frame",
	Run:   cliFrameUpd,
}

func init() {
	rootCmd.AddCommand(cliFrameUpdCmd)
	cliFrameUpdCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
}

func cliFrameUpd(cmd *cobra.Command, args []string) {

	if len(args) != 3 {
		cmd.Help()
		return
	}
	frameId := args[0]
	frameName := args[1]
	frameDomain := args[2]
	payload := `{
			"name": "` + frameName + `",
			"domain": "` + frameDomain + `",
			"webhook": ""
	}`

	apiEndpointPath := "frames/"
	endpoint, _ := cmd.Flags().GetString("path")
	method := "POST"

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, frameId, []byte(payload), "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	var data models.Frame

	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	fmt.Printf("%04d %-32s %45s %s\n", int(data.Id), data.Name, data.Webhook, data.Domain)

}
