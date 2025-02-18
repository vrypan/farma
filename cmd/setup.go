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

	fmt.Println("Generating keys")
	pubKeyHex := config.GetString("key.public")
	privKeyHex := config.GetString("key.private")
	if pubKeyHex != "" {
		fmt.Println(" > A key already exists. Not generarting a new one.")
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
		//viper.Set("key.private", privKeyHex)
		fmt.Println(" >>> Make sure you save your private key somewhere safe! <<<")
		viper.WriteConfig()
	}
	fmt.Printf(" > View/Edit your keypair in %s/%s\n", configDir, "farma.yaml")
	fmt.Println()

	dbPath := db.Path()
	fmt.Printf("Database path is %s\n", dbPath)
	fmt.Println()

	fmt.Printf("The config file is %s/%s\n", configDir, "farma.yaml")
	fmt.Println(" > You can edit it to change these values or use `fargo config`.")
	fmt.Println()
}
