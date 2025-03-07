package cmd

import (
	"encoding/base64"
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
	a := api.ApiCallData{}
	a.Path = "frames/"
	if len(args) != 0 {
		a.Path += args[0]
	}
	a.Endpoint, _ = cmd.Flags().GetString("path")
	a.Method = "POST"

	if len(args) != 2 {
		cmd.Help()
		return
	}
	frameName := args[0]
	frameDomain := args[1]
	a.Body = `{
			"name": "` + frameName + `",
			"domain": "` + frameDomain + `",
			"webhook": ""
	}`

	res, err := a.Call()
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	var data struct {
		Frame   models.Frame
		PrivKey string `json:"privKey"`
	}

	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	fmt.Printf("%04d %-32s %45s %s\n", int(data.Frame.Id), data.Frame.Name, data.Frame.Webhook, data.Frame.Domain)
	encodedPubKey := base64.StdEncoding.EncodeToString([]byte(data.Frame.PublicKey))
	fmt.Printf("Public key: %s\n", encodedPubKey)
	fmt.Printf("Private key: %s\n", data.PrivKey)
	fmt.Println("SAVE YOUR PRIVATE KEY!")
}
