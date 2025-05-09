package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewNotification(
	frameId string,
	appId uint64,
	id string,
	title string,
	message string,
	link string,
	endpoint string,
	tokens map[string]uint64,
) *Notification {
	if id == "" {
		id = uuid.New().String()
	}
	return &Notification{
		FrameId:  frameId,
		AppId:    appId,
		Id:       id,
		Endpoint: endpoint,
		Title:    title,
		Message:  message,
		Link:     link,
		Tokens:   tokens,
	}
}

func (n *Notification) Send_original() error {
	_, err := n.Save()
	if err != nil {
		return fmt.Errorf("Error saving notification: %w", err)
	}

	var tokenKeys []string
	for k := range n.Tokens {
		tokenKeys = append(tokenKeys, k)
	}
	data := map[string]any{
		"notificationId": n.Id,
		"title":          n.Title,
		"body":           n.Message,
		"targetUrl":      n.Link,
		"tokens":         tokenKeys,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshalling json: %w", err)
	}
	request, err := http.NewRequest("POST", n.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating new request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send notification, status code: %d, URL: %s, tokens: %v",
			response.StatusCode, n.Endpoint, n.Tokens)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	// Parse response body and record the status of each token
	var responseBody struct {
		Result struct {
			SuccessfulTokens  []string `json:"successfulTokens"`
			InvalidTokens     []string `json:"invalidTokens"`
			RateLimitedTokens []string `json:"rateLimitedTokens"`
		} `json:"result"`
	}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return fmt.Errorf("Error unmarshalling response body: %w", err)
	}
	n.SuccessTokens = responseBody.Result.SuccessfulTokens
	n.FailedTokens = responseBody.Result.InvalidTokens
	n.RateLimitedTokens = responseBody.Result.RateLimitedTokens
	context := EventContextNotification{Id: n.Id, Version: *n.Version}
	for _, token := range n.SuccessTokens {
		subscription := NewSubscription().FromKey(n.FrameId, n.Tokens[token], n.AppId)
		l := UserLog{
			FrameId:    subscription.FrameId,
			UserId:     subscription.UserId,
			AppId:      subscription.AppId,
			EvtType:    EventType_NOTIFICATION_SENT,
			EvtContext: &UserLog_EventContextNotification{EventContextNotification: &context},
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
	}
	for _, token := range n.FailedTokens {
		subscription := NewSubscription().FromKey(n.FrameId, n.Tokens[token], n.AppId)
		l := UserLog{
			FrameId:    subscription.FrameId,
			UserId:     subscription.UserId,
			AppId:      subscription.AppId,
			EvtType:    EventType_NOTIFICATION_FAILED_INVALID,
			EvtContext: &UserLog_EventContextNotification{EventContextNotification: &context},
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
		subscription.Status = SubscriptionStatus_UNSUBSCRIBED
		subscription.Token = ""
		subscription.Save()
	}
	for _, token := range n.RateLimitedTokens {
		subscription := NewSubscription().FromKey(n.FrameId, n.Tokens[token], n.AppId)
		l := UserLog{
			FrameId:    subscription.FrameId,
			UserId:     subscription.UserId,
			AppId:      subscription.AppId,
			EvtType:    EventType_NOTIFICATION_FAILED_RATE_LIMIT,
			EvtContext: &UserLog_EventContextNotification{EventContextNotification: &context},
		}
		err = l.Save()
		if err != nil {
			return fmt.Errorf("Error in UserLog.Save(): %v\n", err)
		}
		subscription.Status = SubscriptionStatus_RATE_LIMITED
		subscription.Token = token
		subscription.Save()
	}
	return nil
}

func (n *Notification) Prefix() string {
	return "n:id:" + n.FrameId + ":" + n.Id + ":"
}
func (n *Notification) PrefixBytes() []byte {
	return []byte(n.Prefix())
}

func (n *Notification) Save() (int, error) {
	n.Ctime = timestamppb.Now()
	nextVersion := uint64(0)
	var err error
	prefix := n.PrefixBytes()
	for {
		keys, next, err := db.GetKeysWithPrefix(prefix, prefix, 100)
		if err != nil {
			return 0, fmt.Errorf("Error getting keys: %v", err)
		}
		nextVersion += uint64(len(keys))
		if next == nil {
			break
		}
	}
	nextKey := []byte(n.Prefix() + fmt.Sprintf("%03d", nextVersion))
	n.Version = &nextVersion
	notificationBytes, err := proto.Marshal(n)
	if err != nil {
		return 0, fmt.Errorf("Error marshaling notification: %v", err)
	}
	err = db.Set(nextKey, notificationBytes)
	if err != nil {
		return 0, fmt.Errorf("Error saving notification: %v", err)
	}
	return int(nextVersion), nil
}

func (n *Notification) Update() (int, error) {
	key := []byte(n.Prefix() + fmt.Sprintf("%03d", n.GetVersion()))
	notificationBytes, err := proto.Marshal(n)
	if err != nil {
		return 0, fmt.Errorf("Error marshaling notification: %v", err)
	}
	err = db.Set(key, notificationBytes)
	if err != nil {
		return 0, fmt.Errorf("Error saving notification: %v", err)
	}
	return int(*n.Version), nil
}

func (n *Notification) Load(frameId string, notificationId string) ([]*Notification, error) {
	var notifications []*Notification
	n.FrameId = frameId
	n.Id = notificationId
	prefix := n.PrefixBytes()
	next := prefix
	for {
		keys, next, err := db.GetKeysWithPrefix(prefix, next, 100)
		if err != nil {
			return nil, fmt.Errorf("Error getting keys: %v", err)
		}
		if len(keys) > 0 {
			for _, key := range keys {
				value, err := db.Get(key)
				if err != nil {
					return nil, fmt.Errorf("Error getting key %s: %v", key, err)
				}
				notification := Notification{}
				err = proto.Unmarshal(value, &notification)
				if err != nil {
					return nil, fmt.Errorf("Error unmarshaling value for key %s: %v", key, err)
				}
				notifications = append(notifications, &notification)
			}
		}
		if next == nil {
			break
		}
	}
	return notifications, nil
}
