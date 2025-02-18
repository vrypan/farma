package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
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
	cliFramesListCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
}

func cliFramesList(cmd *cobra.Command, args []string) {
	payload := `{
		"command": "frames/get",
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

			id, _ := d["id"].(float64)
			name, _ := d["name"].(string)
			webhook, _ := d["webhook"].(string)
			domain, _ := d["domain"].(string)

			fmt.Printf("%04d %32s %45s %s\n", int(id), name, webhook, domain)
		}
	case "ERROR":
		message, _ := j["message"].(string)
		fmt.Printf("An error occurred: %s\n", message)
	default:
		fmt.Printf("Unknown status received: %s\n", status)
	}
}
