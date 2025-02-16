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

	fmt.Println()
	configDir, err := config.ConfigDir()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Config file is %s/%s\n", configDir, "farma.yaml")
	fmt.Println(" > Make sure you edit it to set your hub.")
	fmt.Println()

	fmt.Println("Generating keys")
	pubKeyHex := config.GetString("key.public")
	privKeyHex := config.GetString("key.private")
	if pubKeyHex != "" || privKeyHex != "" {
		fmt.Println(" > A key already exists. Not generarting a new one.")
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
	fmt.Printf(" > View/Edit your keypair in %s/%s\n", configDir, "farma.yaml")
	fmt.Println()

	dbPath := db.Path()
	fmt.Printf("Database path is %s\n", dbPath)
	fmt.Println()

}
