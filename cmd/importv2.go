package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
)

var importv2Cmd = &cobra.Command{
	Use:   "rawimport <file>",
	Short: "Import data from file",
	Run:   importDatav2,
}

func init() {
	rootCmd.AddCommand(importv2Cmd)
}

func importDatav2(cmd *cobra.Command, args []string) {
	err := db.Open()
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		fmt.Println()
		fmt.Println("If you are running the farma server, you will have to")
		fmt.Println("shut it down in order to use import/export commands.")
		return
	}
	defer db.Close()

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		parts := bytes.SplitN(line, []byte(" "), 2)
		if len(parts) != 2 {
			fmt.Printf("Error parsing line: %s\n", line)
			os.Exit(1)
		}
		key := string(parts[0])
		value := string(parts[1])
		valueBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("Error decoding value for key %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Printf("Setting %s ", key)
		existingValue, err := db.Get([]byte(key))
		if err == nil && existingValue != nil {
			fmt.Println("...skipping")
			continue
		}

		err = db.Set([]byte(key), valueBytes)
		if err != nil {
			fmt.Printf("Error seting key: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Println("OK")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
	}

}
