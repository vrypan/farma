package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
)

type NotificationRequest struct {
	id        string
	frameId   int
	appId     int
	url       string
	tokens    map[string]int
	title     string
	body      string
	targetUrl string
}

func (n NotificationRequest) tokenList() []string {
	ret := make([]string, len(n.tokens))
	index := 0
	for t := range n.tokens {
		ret[index] = t
		index++
	}
	return ret
}

func NewNotificationRequest(title, body, url, targetUrl string, frameId, appId int) *NotificationRequest {
	return &NotificationRequest{
		id:        uuid.New().String(),
		url:       url,
		title:     title,
		body:      body,
		targetUrl: targetUrl,
		frameId:   frameId,
		appId:     appId,
		tokens:    make(map[string]int),
	}
}

func (n *NotificationRequest) AddToken(userId int, token string) {
	n.tokens[token] = userId
}

func (n NotificationRequest) tokenBatches(batchSize int) [][]string {
	allTokens := n.tokenList()
	totalBatches := (len(allTokens) + batchSize - 1) / batchSize
	batches := make([][]string, totalBatches)

	for i := 0; i < len(allTokens); i += batchSize {
		end := i + batchSize
		if end > len(allTokens) {
			end = len(allTokens)
		}
		batches[i/batchSize] = allTokens[i:end]
	}
	return batches
}

func (n NotificationRequest) Send() error {
	// Each notification request sent must have up to 100 tokens.
	// We break them in batches of 100.
	batches := n.tokenBatches(100)

	for _, tokens := range batches {
		data := map[string]interface{}{
			"notificationId": n.id,
			"title":          n.title,
			"body":           n.body,
			"targetUrl":      n.targetUrl + "/" + n.id,
			"tokens":         tokens,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Error marshalling json: %w", err)
		}

		request, err := http.NewRequest("POST", n.url, bytes.NewBuffer(jsonData))
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
				response.StatusCode, n.url, tokens)
		}

		fmt.Println(string(jsonData))
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

		if err == nil {
			statusMap := make(map[string]string)
			tokenStatuses := map[string][]string{
				"Successful":  responseBody.Result.SuccessfulTokens,
				"Invalid":     responseBody.Result.InvalidTokens,
				"RateLimited": responseBody.Result.RateLimitedTokens,
			}

			for status, tokens := range tokenStatuses {
				for _, token := range tokens {
					db.LogUserHistory(
						n.tokens[token],
						n.frameId,
						n.appId,
						fmt.Sprintf("NOTIFICATION_%s", strings.ToUpper(status)),
						n.id,
					)
					statusMap[token] = status
				}
			}
			err = db.UpdateInvalidTokens(tokenStatuses["Invalid"])
			if err != nil {
				return err
			}

			// TO DO: In the case of "RateLimited", users_frames status must be updated to 3

		} else {
			return fmt.Errorf("Error unmarshalling response body: %w", err)
		}
	}
	return nil
}
