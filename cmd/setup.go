package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Printf("%v\n", err)
		os.Exit(1)
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
	fmt.Println()

	fmt.Println("Generating keys")
	pubKeyHex := config.GetString("key.public")
	privKeyHex := config.GetString("key.private")
	if pubKeyHex != "" || privKeyHex != "" {
		fmt.Println(" > The config file already has a key. Not generarting new one.")
		fmt.Printf(" > Private key: %s\n", privKeyHex)
		fmt.Printf(" > Public key: %s\n", pubKeyHex)
	} else {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Printf(" > Error generating keypair: %v\n", err)
			os.Exit(1)
		}
		pubKeyHex = fmt.Sprintf("0x%x", pubKey)
		privKeyHex = fmt.Sprintf("0x%x", privKey)
		fmt.Printf(" > Private key: %s\n", privKeyHex)
		fmt.Printf(" > Public key: %s\n", pubKeyHex)
		viper.Set("key.public", pubKeyHex)
		viper.Set("key.private", privKeyHex)
		viper.WriteConfig()
	}
}
