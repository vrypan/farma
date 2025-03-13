package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vrypan/farma/config"
)

var configgetCmd = &cobra.Command{
	Use:   "get [parameter]",
	Short: "Get config parameter",
	Long: `Examples:
farma config get hub.host

"farma config get" without parameters will return the path
of the configuration file.`,
	Run: config_get,
}

func config_get(cmd *cobra.Command, args []string) {
	f := config.Load()
	if len(args) == 0 {
		fmt.Printf("%s\n", f)
	}
	for _, arg := range args {
		fmt.Printf("%s\n", viper.GetString(arg))
	}
}
func init() {
	configCmd.AddCommand(configgetCmd)
}
