package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var cliFrameAddCmd = &cobra.Command{
	Use:   "frame-add [name] [domain]",
	Short: "Configure a new frame",
	Run:   cliFrameAdd,
}

func init() {
	rootCmd.AddCommand(cliFrameAddCmd)
	cliFrameAddCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
	cliFrameAddCmd.Flags().BoolP("send", "s", true, "Send generated JSON to server")
}

func cliFrameAdd(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Help()
		return
	}
	frameName := args[0]
	frameDomain := args[1]

	payload := `{
		"command": "frames/add",
		"params": {
			"name": "` + frameName + `",
			"domain": "` + frameDomain + `",
			"webhook": ""
		}
	}`
	_, resp := SendCommand(cmd, payload)
	if resp != nil {
		var j map[string]any
		json.Unmarshal(resp, &j)
		if j["status"].(string) == "SUCCESS" {
			data := j["data"].(map[string]any)
			fmt.Printf("Frame Id: %v\n", data["id"].(float64))
			fmt.Printf("Frame Name: %s\n", data["name"].(string))
			fmt.Printf("Frame Domain: %s\n", data["domain"].(string))
			fmt.Printf("Frame Webhook: %s\n", data["webhook"].(string))
		}
		if j["status"].(string) == "ERROR" {
			fmt.Printf("An error occured: %s\n", j["message"].(string))
		}
	}

}
