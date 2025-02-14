package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/config"
)

var rootCmd = &cobra.Command{
	Use:   "farma",
	Short: "A Farcaster notifications manager",
}

func Execute() {
	config.Load()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
