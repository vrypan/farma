package models

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vrypan/farma/fctools"
	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSubscription() *Subscription {
	return &Subscription{}
}
func (s *Subscription) NiceString() string {
	return fmt.Sprintf(
		"FrameId=%s UserId=%d AppId=%d Status=%s Url=%s Token=%s Ctime=%s Mtime=%s AppKey=%s Verified=%t",
		s.FrameId,
		s.UserId,
		s.AppId,
		s.Status.String(),
		s.Url,
		s.Token,
		s.Ctime.AsTime().Format(time.RFC3339),
		s.Mtime.AsTime().Format(time.RFC3339),
		BytesToHex(s.AppKey),
		s.Verified,
	)
}
func (s *Subscription) Key(frameId string, userId, appId uint64) string {
	return fmt.Sprintf("s:id:%s:%d:%d", frameId, userId, appId)
}
func DecodeKey(key []byte) *Subscription {
	s := &Subscription{}
	parts := strings.Split(string(key), ":")
	if len(parts) == 5 {
		frameId := parts[2]
		userId, _ := strconv.ParseUint(parts[3], 10, 64)
		appId, _ := strconv.ParseUint(parts[4], 10, 64)
		s.FrameId = frameId
		s.UserId = userId
		s.AppId = appId
	}
	return s
}
func (s *Subscription) FromHttpEvent(data []byte) (*Subscription, EventType) {
	var jsonBody map[string]any
	if err := json.Unmarshal(data, &jsonBody); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return s, EventType_NONE
	}

	Signature, _ := base64.RawURLEncoding.DecodeString(jsonBody["signature"].(string))
	header, _ := base64.RawURLEncoding.DecodeString(jsonBody["header"].(string))
	payloadDecoded, _ := base64.RawURLEncoding.DecodeString(jsonBody["payload"].(string))

	var headerData map[string]any

	if err := json.Unmarshal(header, &headerData); err == nil {
		s.UserId = uint64(headerData["fid"].(float64))
		s.AppKey = HexToBytes(headerData["key"].(string))
	}

	var payloadData map[string]any
	evtType := EventType_NONE
	if err := json.Unmarshal(payloadDecoded, &payloadData); err == nil {
		eventName := payloadData["event"].(string)
		switch eventName {
		case "frame_added":
			s.Status = SubscriptionStatus_SUBSCRIBED
			if notifDetails, ok := payloadData["notificationDetails"].(map[string]any); ok {
				s.Url = notifDetails["url"].(string)
				s.Token = notifDetails["token"].(string)
			}
		case "notifications_enabled":
			s.Status = SubscriptionStatus_SUBSCRIBED
			if notifDetails, ok := payloadData["notificationDetails"].(map[string]any); ok {
				s.Url = notifDetails["url"].(string)
				s.Token = notifDetails["token"].(string)
			}
		case "frame_removed":
			s.Status = SubscriptionStatus_UNSUBSCRIBED
		case "notifications_disabled":
			s.Status = SubscriptionStatus_UNSUBSCRIBED
		}

		eventNameToType := map[string]EventType{
			"frame_added":            EventType_FRAME_ADDED,
			"notifications_enabled":  EventType_NOTIFICATIONS_ENABLED,
			"frame_removed":          EventType_FRAME_REMOVED,
			"notifications_disabled": EventType_NOTIFICATIONS_DISABLED,
		}
		if eventType, ok := eventNameToType[eventName]; ok {
			evtType = eventType
		}
	}

	signed := []byte(jsonBody["header"].(string) + "." + jsonBody["payload"].(string))
	s.Verified = ed25519.Verify(s.AppKey, signed, Signature)
	return s, evtType
}

func (s *Subscription) VerifyAppId(hub *fctools.FarcasterHub) *Subscription {
	s.AppId = fctools.AppIdFromFidSigner(hub, s.UserId, s.AppKey)
	if s.AppId == 0 {
		s.Verified = false
	}
	return s
}

func (s *Subscription) Save() error {
	subscriptionKey := s.Key(s.FrameId, s.UserId, s.AppId)

	exSub := NewSubscription().FromKey(s.FrameId, s.UserId, s.AppId)
	// Is there an existing subscription in the database?
	if exSub != nil {
		s.Ctime = exSub.GetCtime()
		if exSub.GetToken() != "" {
			tokenKey := NewTokenKey(exSub.GetToken())
			tokenKey.Delete()
		}
	} else {
		s.Ctime = timestamppb.Now()
	}

	tokenKey := NewTokenKey(s.Token)
	s.Mtime = timestamppb.Now()

	data, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	if err = db.Set([]byte(subscriptionKey), data); err != nil {
		return err
	}
	if err = tokenKey.Set([]byte(subscriptionKey)); err != nil {
		return err
	}
	return nil
}

func (s *Subscription) FromKey(frameId string, userId, appId uint64) *Subscription {
	db.AssertOpen()
	key := s.Key(frameId, userId, appId)
	data, err := db.Get([]byte(key))
	if err != nil {
		return nil
	}
	if proto.Unmarshal(data, s) != nil {
		return nil
	}
	return s
}
func (s *Subscription) FromKeyBytes(key []byte) *Subscription {
	data, err := db.Get(key)
	if err != nil {
		return nil
	}
	if proto.Unmarshal(data, s) != nil {
		return nil
	}
	return s
}

func SubscriptionsByFrame(frameId string, start []byte, limit int) ([]*Subscription, []byte, error) {
	prefix := fmt.Appendf([]byte(""), "s:id:%s:", frameId)
	if start == nil {
		start = prefix
	}
	data, nextKey, err := db.GetPrefixP(prefix, start, limit)
	if err != nil {
		return nil, nil, err
	}
	subscriptions := make([]*Subscription, len(data))
	for i, subscription := range data {
		s := NewSubscription()
		if err := proto.Unmarshal(subscription, s); err != nil {
			return nil, nil, err
		}
		subscriptions[i] = s
	}
	return subscriptions, nextKey, nil
}

// Fetch all frame,user subscriptions (all clients). Up to 1000 results.
func SubscriptionsByFrameUser(frameId string, userId uint64) ([]*Subscription, error) {
	prefix := fmt.Appendf([]byte(""), "s:id:%s:%d:", frameId, userId)
	data, _, err := db.GetPrefixP(prefix, prefix, 1000)
	if err != nil {
		return nil, err
	}
	subscriptions := make([]*Subscription, len(data))
	for i, subscription := range data {
		s := NewSubscription()
		if err := proto.Unmarshal(subscription, s); err != nil {
			return nil, err
		}
		subscriptions[i] = s
	}
	return subscriptions, nil
}
