package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/models"
)

var cliFrameAddCmd = &cobra.Command{
	Use:   "frame-add [name] [domain]",
	Short: "Configure a new frame",
	Run:   cliFrameAdd,
}

func init() {
	rootCmd.AddCommand(cliFrameAddCmd)
	cliFrameAddCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
}

func cliFrameAdd(cmd *cobra.Command, args []string) {
	apiEndpointPath := "frames/"
	endpoint, _ := cmd.Flags().GetString("path")

	method := "POST"

	if len(args) != 2 {
		cmd.Help()
		return
	}
	frameName := args[0]
	frameDomain := args[1]
	payload := `{
			"name": "` + frameName + `",
			"domain": "` + frameDomain + `",
			"webhook": ""
	}`

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, "", []byte(payload))
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
