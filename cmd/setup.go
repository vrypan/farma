package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

	dbPath := db.GetDbPath()
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("Database file '%s' already exists. Skipping setup.\n", dbPath)
		return
	}

	fmt.Printf("Creating database at '%s'...\n", dbPath)
	db.Open()
	defer db.Close()
	err := db.CreateTables()
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database setup complete.")
}
