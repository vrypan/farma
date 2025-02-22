package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
)

var cliDbKeysCmd = &cobra.Command{
	Use:   "dbkeys-list [prefix]",
	Short: "List database keys",
	Run:   cliDbKeys,
}

func init() {
	rootCmd.AddCommand(cliDbKeysCmd)
	cliDbKeysCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")

}

func cliDbKeys(cmd *cobra.Command, args []string) {
	apiEndpointPath := "dbkeys/"
	endpoint, _ := cmd.Flags().GetString("path")
	body := []byte("")
	method := "GET"

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	res, err := api.ApiCall(method, endpoint, apiEndpointPath, prefix, body)
	if err != nil {
		fmt.Printf("Failed to make API call: %v", err)
		return
	}

	var list []string
	if err := json.Unmarshal(res, &list); err != nil {
		fmt.Printf("Failed to parse response: %v", err)
		return
	}
	for _, item := range list {
		fmt.Println(item)
	}
}
