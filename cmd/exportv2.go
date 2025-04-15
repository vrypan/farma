package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
)

var cliExportv2Cmd = &cobra.Command{
	Use:   "rawexport",
	Short: "Export db",
	Run:   cliExportv2,
}

func init() {
	rootCmd.AddCommand(cliExportv2Cmd)
	cliExportv2Cmd.Flags().String("path", "", "API endpoint. Defaults to host.addr/api/v2/ (host.addr from config)")
	cliExportv2Cmd.Flags().String("start", "", "Start key")
	cliExportv2Cmd.Flags().String("key", "config", "Private key to use")
}
func cliExportv2(cmd *cobra.Command, args []string) {
	err := db.Open()
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		fmt.Println()
		fmt.Println("If you are running the farma server, you will have to")
		fmt.Println("shut it down in order to use import/export commands.")
		return
	}
	defer db.Close()
	prefix := []byte("")
	next := []byte("")
	for {
		keys, nextKey, err := db.GetKeysWithPrefix(prefix, next, 1000)
		if err != nil {
			fmt.Printf("Error getting keys: %v\n", err)
			return
		}

		for _, key := range keys {
			value, err := db.Get(key)
			if err != nil {
				fmt.Printf("Error getting key %s: %v\n", key, err)
				return
			}
			value64 := base64.StdEncoding.EncodeToString(value)
			fmt.Println(string(key), value64)
		}
		if nextKey == nil {
			break
		}
		next = nextKey
	}

}
