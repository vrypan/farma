package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliFrameUpdCmd = &cobra.Command{
	Use:   "frame-update [frameId]",
	Short: "Update a frame",
	Run:   cliFrameUpd,
}

func init() {
	rootCmd.AddCommand(cliFrameUpdCmd)
	cliFrameUpdCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliFrameUpdCmd.Flags().String("public-key", "", "Base64 encoded public key")
	cliFrameUpdCmd.Flags().String("name", "", "Frame name")
	cliFrameUpdCmd.Flags().String("domain", "", "Frame domain")
}

func cliFrameUpd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}
	frameId := args[0]
	name, _ := cmd.Flags().GetString("name")
	domain, _ := cmd.Flags().GetString("domain")
	public_key, _ := cmd.Flags().GetString("public-key")

	payload := `{
		"name": "` + name + `",
		"domain": "` + domain + `",
		"public_key": "` + public_key + `"
		}`
	a := api.ApiClient{}.Init("POST", "frame/"+frameId, []byte(payload), []byte("config"), "")
	res, err := a.Request("", "")
	if err != nil {
		fmt.Printf("Failed to make API call: %v %s\n", err, res)
		return
	}
	fmt.Println(string(res))

	fmt.Printf("MAKE SURE YOU SAVE THE ABOVE INFORMATION!")
	fmt.Printf("Many API calls, require the FrameID, the PublicKey, and the PrivateKey.")

}
