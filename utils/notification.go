package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
)

func NewNotification(title, message string, link string, endpoint string, urlKeys [][]byte) *Notification {
	tokens := make([]string, len(urlKeys))
	urlKey := UrlKey{}
	for i, key := range urlKeys {
		tokens[i] = urlKey.DecodeBytes(key).Token
	}
	return &Notification{
		Id:       uuid.New().String(),
		Endpoint: endpoint,
		Title:    title,
		Message:  message,
		Link:     link,
		Tokens:   tokens,
	}
}

func (n *Notification) Send() error {
	data := map[string]interface{}{
		"notificationId": n.Id,
		"title":          n.Title,
		"body":           n.Message,
		"targetUrl":      n.Link,
		"tokens":         n.Tokens,
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
	for _, token := range n.SuccessTokens {
		tokenKey := NewTokenKey(token)
		subscriptionKey, err := db.Get(tokenKey.Bytes())
		if err != nil {
			return fmt.Errorf("Error getting subscription key: %w", err)
		}
		if subscriptionKey == nil {
			return fmt.Errorf("Subscription key not found for token: %s", token)
		}
		subscription := DecodeKey(subscriptionKey)
		l := UserLog{
			FrameId:      subscription.FrameId,
			UserId:       subscription.UserId,
			AppId:        subscription.AppId,
			EvtType:      EventType_NOTIFICATION_SENT,
			EventContext: token,
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
	}
	for _, token := range n.FailedTokens {
		tokenKey := NewTokenKey(token)
		subscriptionKey, err := db.Get(tokenKey.Bytes())
		if err != nil {
			return fmt.Errorf("Error getting subscription key: %v", err)
		}
		if subscriptionKey == nil {
			return fmt.Errorf("Subscription key not found for token: %s", token)
		}
		subscription := DecodeKey(subscriptionKey)
		l := UserLog{
			FrameId:      subscription.FrameId,
			UserId:       subscription.UserId,
			AppId:        subscription.AppId,
			EvtType:      EventType_NOTIFICATION_FAILED_INVALID,
			EventContext: token,
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
		subscription.FromKeyBytes(subscriptionKey)
		subscription.Status = SubscriptionStatus_UNSUBSCRIBED
		subscription.Token = ""
		subscription.Save()
	}
	for _, token := range n.RateLimitedTokens {
		tokenKey := NewTokenKey(token)
		subscriptionKey, err := db.Get(tokenKey.Bytes())
		if err != nil {
			return fmt.Errorf("Error getting subscription key: %v", err)
		}
		if subscriptionKey == nil {
			return fmt.Errorf("Subscription key not found for token: %s", token)
		}
		subscription := DecodeKey(subscriptionKey)
		l := UserLog{
			FrameId:      subscription.FrameId,
			UserId:       subscription.UserId,
			AppId:        subscription.AppId,
			EvtType:      EventType_NOTIFICATION_FAILED_RATE_LIMIT,
			EventContext: token,
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
		subscription.FromKeyBytes(subscriptionKey)
		subscription.Status = SubscriptionStatus_RATE_LIMITED
		subscription.Save()
	}
	return nil
}
