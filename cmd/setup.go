package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
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
	Long: `Setup will create a new config file and an admin keypair.
By default, the config directory is ~/.farma/, but
XDG_CONFIG_HOME is honored if set. After the setup,
you can use "farma config" to get/set values.`,
	Run: setupFidr,
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
	pubKey64 := config.GetString("key.public")
	privKey64 := config.GetString("key.private")
	if pubKey64 != "" {
		fmt.Println(" > A key already exists. Not generarting a new one.")
		fmt.Printf(" > Public key: %s\n", pubKey64)
	} else {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Printf(" > Error generating keypair: %v\n", err)
			os.Exit(1)
		}
		//pubKeyHex = fmt.Sprintf("0x%x", pubKey)
		//privKeyHex = fmt.Sprintf("0x%x", privKey)
		pubKey64 = base64.StdEncoding.EncodeToString(pubKey)
		privKey64 = base64.StdEncoding.EncodeToString(privKey)
		fmt.Printf(" > Private key: %s\n", privKey64)
		fmt.Printf(" > Public key: %s\n", pubKey64)
		viper.Set("key.public", pubKey64)
		viper.Set("key.private", privKey64)
		fmt.Printf(" >>> Keys are stored in %s/config.yaml\n", configDir)
		viper.WriteConfig()
	}
	fmt.Printf(" > View/Edit your keypair in %s/%s\n", configDir, "config.yaml")
	fmt.Println()

	dbPath := db.Path()
	fmt.Printf("Database path is %s\n", dbPath)
	fmt.Println()

	fmt.Printf("The config file is %s/%s\n", configDir, "config.yaml")
	fmt.Println(" > You can edit it to change these values or use `fargo config`.")
	fmt.Println()
}
