package cmd

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
)

var frameAddCmd = &cobra.Command{
	Use:     "add [NAME]",
	Short:   "Add a new frame",
	Aliases: []string{"a"},
	Run:     addFrame,
}

func init() {
	frameCmd.AddCommand(frameAddCmd)
	frameAddCmd.Flags().StringP("description", "", "", "Optional long description")
	frameAddCmd.Flags().StringP("webhook", "", "", "Endpoint path") // webhook is automatically generated
	frameAddCmd.Flags().StringP("domain", "", "", "Frame domain")
}

func addFrame(cmd *cobra.Command, args []string) {
	domain, _ := cmd.Flags().GetString("domain")
	webhook, _ := cmd.Flags().GetString("webhook")

	db.Open()
	defer db.Close()

	if len(args) == 0 {
		fmt.Fprintln(cmd.OutOrStderr(), "Error: Frame NAME not defined")
		os.Exit(1)
	}
	name := args[0]
	var endpoint string
	if webhook == "" {
		endpoint = "/f/" + uuid.New().String()
	} else {
		endpoint = webhook
	}

	if len(name) > 32 {
		fmt.Fprintln(cmd.OutOrStderr(), "Name must be up to 32 characters")
		os.Exit(1)
	}
	//sql := fmt.Sprintf("INSERT INTO frames (name, domain, endpoint) VALUES('%s','/frame/%s')", name, uuid)
	_, err := db.Instance.Exec("INSERT INTO frames (name, domain, endpoint) VALUES(?,?,?)",
		name, domain, endpoint,
	)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error:", err)
		os.Exit(1)
	}
	fmt.Println("Frame added")
}
