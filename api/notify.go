package api

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

// Send out notifications.
func Notify(frame string, title string, body string, url string) error {
	db.AssertOpen()

	rows, err := db.Instance.Query(`
		SELECT id, name, desc, domain, endpoint FROM frames WHERE name=?
		`, frame)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return fmt.Errorf("Frame %s not found", frame)
	}

	var frameId int
	var frameName string
	var frameDesc string
	var frameDomain string
	var frameEndpoint string

	err = rows.Scan(&frameId, &frameName, &frameDesc, &frameDomain, &frameEndpoint)
	if err != nil {
		return err
	}

	//fmt.Printf("Frame - ID: %04d Name: %-32s Endpoint: %s Description: %s\n", frameId, frameName, frameEndpoint, frameDesc)

	rows, err = db.Instance.Query(`
		SELECT id, user_id, app_id, status, url, token
		FROM users_frames WHERE frame_id=? and status=1
		`, frameId)
	if err != nil {
		return err
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
			return fmt.Errorf("Error scanning row. %v", err)
		}

		_, ok := notificationRequests[appEndpoint]
		if !ok {
			if len(url) > 0 {
				targetUrl = url
			} else {
				targetUrl = "https://" + frameDomain
			}
			notificationRequests[appEndpoint] = utils.NewNotificationRequest(
				title,
				body,
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
	return nil
}
