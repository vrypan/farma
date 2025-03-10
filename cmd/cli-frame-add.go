package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
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
	cliFrameAddCmd.Flags().String("key", "config", "Private key to use.")
}

func cliFrameAdd(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
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
	a := api.ApiClient{}.Init("POST", "/api/v2/frame/", []byte(payload), key, "")
	res, err := a.Request("", "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	fmt.Println(string(res))

	fmt.Println()
	fmt.Println("MAKE SURE YOU SAVE THE FOLLOWING INFORMATION!")
	fmt.Println("Many API calls, require the FrameID, the PublicKey, and the PrivateKey.")

	var data struct {
		Frame      models.Frame `json:"frame"`
		PublicKey  string       `json:"public_key"`
		PrivateKey string       `json:"private_key"`
	}
	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	//var frame models.Frame
	fmt.Println()
	fmt.Println("Frame ID:   ", data.Frame.Id)
	fmt.Println("Public Key: ", data.PublicKey)
	fmt.Println("Private Key:", data.PrivateKey)
	fmt.Println("Endpoint:", "/f/"+data.Frame.Id)
}
