package api

import (
	"strconv"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

// Send out notifications.
func Notify(frameName string, notificationTitle string, notificationBody string, notificationUrl string) string {
	response := Response{}
	frame := utils.NewFrame()
	if frame.FromName(frameName) != nil {
		return response.Format("error", "FRAME_NOT_FOUND", nil)
	}

	// Warpcast will crash when an notificationUrl is clicked.
	if notificationUrl == "" {
		notificationUrl = "https://" + frame.Domain
	}

	keys := make(map[string][][]byte)

	prefix := []byte("s:url:" + strconv.Itoa(int(frame.Id)) + ":")

	startKey := prefix
	for {
		urlKeys, nextKey, err := db.GetKeysWithPrefix(prefix, startKey, 1000)
		if err != nil {
			return Error("DB_ERROR", err)
		}
		for _, urlKeyBytes := range urlKeys {
			urlKey := utils.UrlKey{}.DecodeBytes(urlKeyBytes)
			// urlkey is s:url:<frameId>:<userId>:<status>:<url>:<token>
			status := urlKey.Status
			url := urlKey.Endpoint
			if status == utils.SubscriptionStatus_SUBSCRIBED || status == utils.SubscriptionStatus_RATE_LIMITED {
				keys[url] = append(keys[url], urlKeyBytes)
			}
		}
		startKey = nextKey
		if len(urlKeys) < 1000 {
			break
		}
	}

	for url, urlKeys := range keys {
		notification := utils.NewNotification(
			notificationTitle,
			notificationBody,
			notificationUrl,
			url,
			urlKeys,
		)
		err := notification.Send()
		if err != nil {
			return Error("NOTIFICATION_ERROR", err)
		}
	}
	return response.Format("SUCCESS", "Notification sent", nil)
}
