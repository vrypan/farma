package apiv2

import (
	"github.com/vrypan/farma/models"
)

/*
Warpcast exports data in this format:

	fid,notificationToken,added
	2,0195fbc2-0000-0000-0000-29670ebd1097,true
	3,null,false
	280,019643a2-0000-0000-0000-bf2bd584a469,true
*/
type csvEntry struct {
	fid               uint64
	notificationToken string
	added             bool
}

type ImportData struct {
	appId   uint64 // Application FID
	appUrl  string // Application notifications endpoint
	frameId string
	data    []csvEntry
}

func CreateSubscriptionsFromCSV(d ImportData) (int, error) {
	count := 0
	for _, entry := range d.data {
		status := models.SubscriptionStatus_SUBSCRIBED
		token := entry.notificationToken

		if len(entry.notificationToken) == 0 || entry.notificationToken == "null" {
			status = models.SubscriptionStatus_UNSUBSCRIBED
			token = ""
		}

		subscription := &models.Subscription{
			FrameId: d.frameId,
			UserId:  entry.fid,
			AppId:   d.appId,
			Status:  status,
			Url:     d.appUrl,
			Token:   token,
		}

		// Assume SaveSubscription is a method to save the subscription in the database
		if err := subscription.Save(); err != nil {
			return count, err
		}

		count++
	}
	return count, nil
}
