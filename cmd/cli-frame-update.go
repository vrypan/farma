package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
)

var cliFrameUpdCmd = &cobra.Command{
	Use:   "frame-update [id]",
	Short: "Update a frame",
	Run:   cliFrameUpd,
}

func init() {
	rootCmd.AddCommand(cliFrameUpdCmd)
	cliFrameUpdCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliFrameUpdCmd.Flags().String("public-key", "", "Base64 encoded public key")
	cliFrameUpdCmd.Flags().String("name", "", "Frame name")
	cliFrameUpdCmd.Flags().String("domain", "", "Frame domain")
	cliFrameUpdCmd.Flags().String("webhook", "", "Frame webhook")
}

func cliFrameUpd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}

	a := api.ApiCallData{}
	a.Path = "frames/" + args[0]
	if len(args) != 0 {
		a.Path += args[0]
	}
	a.Endpoint, _ = cmd.Flags().GetString("path")
	a.Method = "POST"

	name, _ := cmd.Flags().GetString("name")
	domain, _ := cmd.Flags().GetString("domain")
	webhook, _ := cmd.Flags().GetString("webhook")
	public_key, _ := cmd.Flags().GetString("public-key")
	a.Body = `{
			"name": "` + name + `",
			"domain": "` + domain + `",
			"webhook": "` + webhook + `",
			"public_key": "` + public_key + `"
	}`

	res, err := a.Call()
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	/*
		var data models.Frame

		if err := json.Unmarshal(res, &data); err != nil {
			fmt.Printf("Failed to parse response: %v", err)
			return
		}
	*/
	fmt.Printf("%s\n", string(res))

}
