package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var requestBody struct {
	FrameId string   `json:"frameId"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Url     string   `json:"url"`
	UserIds []uint64 `json:"userIds"`
}

func PrepareNotification() {

}

func (n *Notification) Send() error {
	_, err := n.Save()
	if err != nil {
		return fmt.Errorf("Error saving notification: %w", err)
	}

	var tokenKeys []string
	for k := range n.Tokens {
		tokenKeys = append(tokenKeys, k)
	}

	const batchSize = 100
	batchNumber := 0
	for i := 0; i < len(tokenKeys); i += batchSize {
		batchNumber += 1
		end := i + batchSize
		end = min(end, len(tokenKeys))

		batchTokens := tokenKeys[i:end]

		data := map[string]any{
			"notificationId": n.Id,
			"title":          n.Title,
			"body":           n.Message,
			"targetUrl":      n.Link,
			"tokens":         batchTokens,
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
			n.ServerErrorTokens = append(n.FailedTokens, batchTokens...)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			n.ServerErrorTokens = append(n.FailedTokens, batchTokens...)
			continue
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
		n.SuccessTokens = append(n.SuccessTokens, responseBody.Result.SuccessfulTokens...)
		n.FailedTokens = append(n.FailedTokens, responseBody.Result.InvalidTokens...)
		n.RateLimitedTokens = append(n.RateLimitedTokens, responseBody.Result.RateLimitedTokens...)
	}

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
	for _, token := range n.ServerErrorTokens {
		subscription := NewSubscription().FromKey(n.FrameId, n.Tokens[token], n.AppId)
		l := UserLog{
			FrameId:    subscription.FrameId,
			UserId:     subscription.UserId,
			AppId:      subscription.AppId,
			EvtType:    EventType_NOTIFICATION_FAILED_OTHER,
			EvtContext: &UserLog_EventContextNotification{EventContextNotification: &context},
		}
		err = l.Save()
		if err != nil {
			fmt.Printf("Error in UserLog.Save(): %v\n", err)
		}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
