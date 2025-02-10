package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
)

var frameLsCmd = &cobra.Command{
	Use:     "list [frame name",
	Short:   "List frames",
	Aliases: []string{"ls"},
	Run:     lsFrame,
}

func init() {
	frameCmd.AddCommand(frameLsCmd)
	frameLsCmd.Flags().BoolP("with-subscriptions", "", false, "Show subscriptions for each frame")
}

func lsFrame(cmd *cobra.Command, args []string) {
	subscriptionsFlag, _ := cmd.Flags().GetBool("with-subscriptions")

	db.Open()
	defer db.Close()

	rows, err := db.Instance.Query("SELECT id, name, domain, endpoint FROM frames ORDER BY id ASC")
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error:", err)
		os.Exit(1)
	}

	for rows.Next() {
		var id int
		var name string
		var domain string
		var endpoint string
		err := rows.Scan(&id, &name, &domain, &endpoint)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "Error scanning row:", err)
			continue
		}

		fmt.Printf("%04d %-32s %s %s\n", id, name, endpoint, domain)
		if subscriptionsFlag {
			listSubscriptions(id)
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
}

func listSubscriptions(frameId int) {
	rows, err := db.Instance.Query(`
		SELECT user_id, app_id, status, token, ctime, mtime
		FROM users_frames
		WHERE frame_id=?
		`, frameId)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for rows.Next() {
		var userId int
		var appId int
		var status bool
		var token string
		var ctime string
		var mtime string
		err := rows.Scan(&userId, &appId, &status, &token, &ctime, &mtime)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		fmt.Printf("> %012d %012d %v %s %s %s\n", userId, appId, status, token, ctime, mtime)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
}
