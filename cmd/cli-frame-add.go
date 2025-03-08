package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
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
	a := api.ApiClient{}.Init("POST", "frame/", []byte(payload), []byte("config"), "")
	res, err := a.Request("", "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	fmt.Println(string(res))

	fmt.Printf("MAKE SURE YOU SAVE THE ABOVE INFORMATION!")
	fmt.Printf("Many API calls, require the FrameID, the PublicKey, and the PrivateKey.")

}
