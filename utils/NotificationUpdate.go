package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vrypan/farma/fctools"
	db "github.com/vrypan/farma/localdb"
)

type EventType int

const (
	FRAME_ADDED EventType = iota
	FRAME_REMOVED
	NOTIFICATIONS_ENABLED
	NOTIFICATIONS_DISABLED
	NOTIFICATION_SENT
	NOTIFICATION_FAILED
)

func (t EventType) String() string {
	eventNames := map[EventType]string{
		FRAME_ADDED:            "FRAME_ADDED",
		FRAME_REMOVED:          "FRAME_REMOVED",
		NOTIFICATIONS_ENABLED:  "NOTIFICATIONS_ENABLED",
		NOTIFICATIONS_DISABLED: "NOTIFICATIONS_DISABLED",
		NOTIFICATION_SENT:      "NOTIFICATION_SENT",
		NOTIFICATION_FAILED:    "NOTIFICATION_FAILED",
	}
	if name, ok := eventNames[t]; ok {
		return name
	}
	return ""
}

// NotificationUpdate is an update sent from a Farcaster client
// (ex. Warpcast) about the subscription status of a frame.
type NotificationUpdate struct {
	FrameId   int
	Type      EventType
	UserFid   int
	AppKey    string
	AppFid    int
	AppUrl    string
	Token     string
	Signature []byte
	Verified  bool
}

func NewNotificationUpdate(frameId int) *NotificationUpdate {
	return &NotificationUpdate{FrameId: frameId}
}

func (n *NotificationUpdate) FromHttpEvent(data []byte) *NotificationUpdate {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(data, &jsonBody); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return n
	}

	n.Signature, _ = base64.RawURLEncoding.DecodeString(jsonBody["signature"].(string))
	header, _ := base64.RawURLEncoding.DecodeString(jsonBody["header"].(string))
	payloadDecoded, _ := base64.RawURLEncoding.DecodeString(jsonBody["payload"].(string))

	var headerData map[string]interface{}
	if err := json.Unmarshal(header, &headerData); err == nil {
		n.UserFid = int(headerData["fid"].(float64))
		n.AppKey = headerData["key"].(string)
	}

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payloadDecoded, &payloadData); err == nil {
		eventName := payloadData["event"].(string)
		eventTypes := map[string]EventType{
			"frame_added":            FRAME_ADDED,
			"notifications_enabled":  NOTIFICATIONS_ENABLED,
			"frame_removed":          FRAME_REMOVED,
			"notifications_disabled": NOTIFICATIONS_DISABLED,
		}
		if eventType, ok := eventTypes[eventName]; ok {
			n.Type = eventType
			if eventName == "frame_added" || eventName == "notifications_enabled" {
				if notifDetails, ok := payloadData["notificationDetails"].(map[string]interface{}); ok {
					n.AppUrl = notifDetails["url"].(string)
					n.Token = notifDetails["token"].(string)
				}
			}
		}
	}

	signed := []byte(jsonBody["header"].(string) + "." + jsonBody["payload"].(string))
	pubKey := common.FromHex(n.AppKey)
	n.Verified = ed25519.Verify(pubKey, signed, n.Signature)
	return n
}

func (n NotificationUpdate) IsActive() int {
	isActiveEvents := map[EventType]int{
		FRAME_ADDED:            1,
		NOTIFICATIONS_ENABLED:  1,
		FRAME_REMOVED:          0,
		NOTIFICATIONS_DISABLED: 0,
	}

	if isActive, ok := isActiveEvents[n.Type]; ok {
		return isActive
	}

	return 0
}

// Get the AppFid associated with the signing key.
// If we fail to get the AppFid, set Verified=false
func (n *NotificationUpdate) GetAppFid(hub *fctools.FarcasterHub) *NotificationUpdate {
	n.AppFid = int(fctools.AppIdFromFidSigner(hub, uint64(n.UserFid), common.FromHex(n.AppKey)))
	if n.AppFid == 0 {
		n.Verified = false
	}
	return n
}

func (s NotificationUpdate) String() string {
	return fmt.Sprintf("%s FrameId=%d UserFid=%d AppFid=%d AppKey=%s AppUrl=%s AppToken=%s",
		s.Type, s.FrameId, s.UserFid, s.AppFid, s.AppKey, s.AppUrl, s.Token)
}

func (n NotificationUpdate) UpdateDb() error {
	db.AssertOpen()
	if n.Verified == false {
		return fmt.Errorf("Not verified. Will not update db: %012d-%012d-%012d", n.UserFid, n.FrameId, n.AppFid)
	}

	err := db.UpdateFrameStatus(n.UserFid, n.FrameId, n.AppFid, (n.IsActive() == 1), n.AppUrl, n.Token)
	if err != nil {
		return err
	}
	err = db.LogUserHistory(n.UserFid, n.FrameId, n.AppFid, n.Type.String(), "")
	if err != nil {
		return err
	}
	return nil
}
