package api

import (
	"fmt"
	"strconv"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

// Send out notifications.
func Notify(frameName string, notificationTitle string, notificationBody string, notificationUrl string) error {

	frame := utils.NewFrame()
	if frame.FromName(frameName) == nil {
		fmt.Println("Frame not found")
		return fmt.Errorf("Frame not found")
	}
	fmt.Println(frame)

	keys := make(map[string][][]byte)

	prefix := []byte("s:url:" + strconv.Itoa(int(frame.GetId())) + ":")
	startKey := prefix
	for {
		urlKeys, nextKey, err := db.GetPrefixP(prefix, startKey, 1000)
		if err != nil {
			return fmt.Errorf("Error fetching subscriptions: %v", err)
		}
		if len(urlKeys) < 1000 {
			break
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
		fmt.Printf("Notification %s sent with result %v\n", notification.Id, err)
	}
	return nil
}
