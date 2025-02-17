package utils

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
		"FrameId=%d UserId=%d AppId=%d Status=%s Url=%s Token=%s Ctime=%s Mtime=%s AppKey=%s Verified=%t",
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
func (s *Subscription) Key(frameId, userId, appId uint64) string {
	return fmt.Sprintf("s:id:%d:%d:%d", frameId, userId, appId)
}
func DecodeKey(key []byte) *Subscription {
	s := &Subscription{}
	parts := strings.Split(string(key), ":")
	if len(parts) == 4 {
		frameId, _ := strconv.ParseUint(parts[1], 10, 64)
		userId, _ := strconv.ParseUint(parts[2], 10, 64)
		appId, _ := strconv.ParseUint(parts[3], 10, 64)
		s.FrameId = frameId
		s.UserId = userId
		s.AppId = appId
	}
	return s
}
func (s *Subscription) FromHttpEvent(data []byte) *Subscription {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(data, &jsonBody); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return s
	}

	Signature, _ := base64.RawURLEncoding.DecodeString(jsonBody["signature"].(string))
	header, _ := base64.RawURLEncoding.DecodeString(jsonBody["header"].(string))
	payloadDecoded, _ := base64.RawURLEncoding.DecodeString(jsonBody["payload"].(string))

	var headerData map[string]interface{}

	if err := json.Unmarshal(header, &headerData); err == nil {
		s.UserId = uint64(headerData["fid"].(float64))
		s.AppKey = HexToBytes(headerData["key"].(string))
	}

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payloadDecoded, &payloadData); err == nil {
		eventName := payloadData["event"].(string)
		switch eventName {
		case "frame_added":
		case "notifications_enabled":
			s.Status = SubscriptionStatus_SUBSCRIBED
			if notifDetails, ok := payloadData["notificationDetails"].(map[string]interface{}); ok {
				s.Url = notifDetails["url"].(string)
				s.Token = notifDetails["token"].(string)
			}
		case "frame_removed":
		case "notifications_disabled":
			s.Status = SubscriptionStatus_UNSUBSCRIBED
		}
	}

	signed := []byte(jsonBody["header"].(string) + "." + jsonBody["payload"].(string))
	s.Verified = ed25519.Verify(s.AppKey, signed, Signature)
	return s
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

	newTokenKey := fmt.Sprintf("s:token:%s", s.Token)
	oldTokenKey := ""

	newUrlKey := UrlKey{}.FromSubscription(s)
	oldUrlKey := UrlKey{}

	if s.Ctime == nil {
		tmp := NewSubscription().FromKey(s.FrameId, s.UserId, s.AppId)
		if tmp != nil {
			// If Subscription was already saved in DB
			s.Ctime = tmp.Ctime // Inherit ctime
			oldTokenKey = fmt.Sprintf("s:token:%s", tmp.Token)
			oldUrlKey = oldUrlKey.FromSubscription(tmp)
		}
	}
	s.Mtime = timestamppb.Now()

	data, err := proto.Marshal(s)
	if err != nil {
		return err
	}

	if err = db.Set([]byte(subscriptionKey), data); err != nil {
		return err
	}
	if oldTokenKey != newTokenKey {
		if oldTokenKey != "" {
			if err := db.Delete([]byte(oldTokenKey)); err != nil {
				return err
			}
			if err := db.Delete(oldUrlKey.Bytes()); err != nil {
				return err
			}
		}
		if newTokenKey != "" {
			if err := db.Set([]byte(newTokenKey), []byte(subscriptionKey)); err != nil {
				return err
			}
			if err := db.Set(newUrlKey.Bytes(), []byte(subscriptionKey)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Subscription) FromKey(frameId, userId, appId uint64) *Subscription {
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

func SubscriptionsByFrame(frameId uint64, limit int) ([]*Subscription, []byte, error) {
	prefix := fmt.Sprintf("s:id:%d:", frameId)
	data, nextKey, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)
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
