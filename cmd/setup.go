package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/config"
	db "github.com/vrypan/farma/localdb"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize and setup farma",
	Run:   setupFidr,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func setupFidr(cmd *cobra.Command, args []string) {

	configDir, err := config.ConfigDir()
	if err != nil {
		fmt.Printf("%w\n", err)
		return
	}
	fmt.Printf("Config file is %s/%s\n", configDir, "farma.yaml")
	fmt.Println(" > Make sure you edit it to set your hub.")
	fmt.Println()

	dbPath := db.GetDbPath()
	fmt.Printf("Database file is %s\n", dbPath)
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println(" > File already exists. Skipping.")
	} else {
		fmt.Println(" > Creating database...")
		db.Open()
		defer db.Close()
		err := db.CreateTables()
		if err != nil {
			fmt.Printf(" > Error creating database: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(" > Database setup complete.")
	}

}
