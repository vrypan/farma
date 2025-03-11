package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/vrypan/farma/apiv2"
)

var cliDbKeysCmd = &cobra.Command{
	Use:   "dbkeys-list [prefix]",
	Short: "List database keys",
	Run:   cliDbKeys,
}

func init() {
	rootCmd.AddCommand(cliDbKeysCmd)
	cliDbKeysCmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v1/frames/ (from config file)")
	cliDbKeysCmd.Flags().String("start", "", "Start key")
	cliDbKeysCmd.Flags().Int("limit", 1000, "Max results")
	cliDbKeysCmd.Flags().String("key", "config", "Private key to use.")
}

func cliDbKeys(cmd *cobra.Command, args []string) {
	key, _ := cmd.Flags().GetString("key")
	start, _ := cmd.Flags().GetString("start")
	limit, _ := cmd.Flags().GetInt("limit")
	path, _ := cmd.Flags().GetString("path")
	path += "dbkeys/"

	if len(args) > 0 {
		path = path + args[0]
	}

	a := api.ApiClient{}.Init("GET", path, nil, key, "")

	next := start
	count := 0
	for {
		if count >= limit {
			break
		}
		resBytes, err := a.Request(next, fmt.Sprintf("%d", limit))
		if err != nil {
			fmt.Printf("Failed to make API call: %v\n", err)
			fmt.Println(string(resBytes))
			return
		}
		var res api.ApiResult
		if err := json.Unmarshal(resBytes, &res); err != nil {
			fmt.Printf("Failed to parse response: %v\n", err)
			return
		}
		for _, v := range res.Result {
			fmt.Println(string(v))
		}
		if res.Next == "" {
			break
		}
		next = res.Next
		count += len(res.Result)
	}
}
