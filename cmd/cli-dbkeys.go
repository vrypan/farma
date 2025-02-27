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
	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}
	endpoint, _ := cmd.Flags().GetString("path")
	body := []byte("")
	method := "GET"

	next := ""
	for {
		resBytes, err := api.ApiCall(method, endpoint, apiEndpointPath, prefix, body, "start="+next)
		if err != nil {
			fmt.Printf("Failed to make API call: %v", err)
			return
		}

		var res struct {
			Error  string   `json:"error"`
			Result []string `json:"result"`
			Next   string   `json:"next"`
		}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Printf("Failed to parse response: %v", err)
			return
		}
		for _, v := range res.Result {
			fmt.Println(string(v))
		}
		if res.Next == "" {
			break
		}
		next = res.Next
	}
}
