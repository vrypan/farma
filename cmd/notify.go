package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var notifyCmd = &cobra.Command{
	Use:   "notify [frame name]",
	Short: "Send a notification to all [frame name] users. If [fid] is provided only send to [fid].",
	Run:   notify,
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().String("title", "", "Notification title")
	notifyCmd.Flags().String("body", "", "Notification body")
	notifyCmd.Flags().String("url", "", "Oprional target URL")
}

func notify(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		fmt.Printf("You must provide a frame name.\n")
		os.Exit(1)
	}
	notificationTitle, _ := cmd.Flags().GetString("title")
	if len(notificationTitle) == 0 {
		fmt.Println("Use --title to set the notification title")
		os.Exit(1)
	}
	notificationBody, _ := cmd.Flags().GetString("body")
	if len(notificationBody) == 0 {
		fmt.Println("Use --body to set the notification body")
		os.Exit(1)
	}
	notificationUrl, _ := cmd.Flags().GetString("url")

	db.Open()
	defer db.Close()

	rows, err := db.Instance.Query(`
		SELECT id, name, desc, domain, endpoint FROM frames WHERE name=?
		`, args[0])
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error:", err)
		os.Exit(1)
	}
	if !rows.Next() {
		fmt.Fprintln(cmd.OutOrStderr(), "Frame not found")
		return
	}

	var frameId int
	var frameName string
	var frameDesc string
	var frameDomain string
	var frameEndpoint string

	err = rows.Scan(&frameId, &frameName, &frameDesc, &frameDomain, &frameEndpoint)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error scanning frame data:", err)
		os.Exit(1)
	}

	fmt.Printf("Frame - ID: %04d Name: %-32s Endpoint: %s Description: %s\n", frameId, frameName, frameEndpoint, frameDesc)

	rows, err = db.Instance.Query(`
		SELECT id, user_id, app_id, status, url, token
		FROM users_frames WHERE frame_id=? and status=1
		`, frameId)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Error:", err)
		os.Exit(1)
	}

	var id string
	var userId int
	var appId int
	var status int
	var appEndpoint string
	var token string
	var targetUrl string

	notificationRequests := make(map[string]*utils.NotificationRequest)

	for rows.Next() {
		err := rows.Scan(&id, &userId, &appId, &status, &appEndpoint, &token)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "Error scanning row:", err)
			continue
		}

		_, ok := notificationRequests[appEndpoint]
		if !ok {
			if len(notificationUrl) > 0 {
				targetUrl = notificationUrl
			} else {
				targetUrl = "https://" + frameDomain
			}
			notificationRequests[appEndpoint] = utils.NewNotificationRequest(
				notificationTitle,
				notificationBody,
				appEndpoint,
				targetUrl,
				frameId,
				appId,
			)
		}
		notificationRequests[appEndpoint].AddToken(userId, token)
	}
	if len(notificationRequests) > 0 {
		for _, n := range notificationRequests {
			n.Send()
		}
	}
}
