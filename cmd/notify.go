package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var notifyCmd = &cobra.Command{
	Use:   "notify [frame name]",
	Short: "Send a notification to all [frame name] users. If [fid] is provided only send to [fid].",
	Run:   notify,
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().String("title", "", "Notification title")
	notifyCmd.Flags().String("body", "", "Notification body")
	notifyCmd.Flags().String("url", "", "Oprional target URL")
}

func notify(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		fmt.Printf("You must provide a frame name.\n")
		os.Exit(1)
	}
	frameName := args[0]
	notificationTitle, _ := cmd.Flags().GetString("title")
	if len(notificationTitle) == 0 {
		fmt.Println("Use --title to set the notification title")
		os.Exit(1)
	}
	notificationBody, _ := cmd.Flags().GetString("body")
	if len(notificationBody) == 0 {
		fmt.Println("Use --body to set the notification body")
		os.Exit(1)
	}
	notificationUrl, _ := cmd.Flags().GetString("url")

	db.Open()
	defer db.Close()

	frame := utils.NewFrame()
	if frame.FromName(frameName) == nil {
		fmt.Fprintln(cmd.OutOrStderr(), "Frame not found")
		return
	}
	fmt.Println(frame)

	keys := make(map[string][][]byte)

	prefix := []byte("s:url:" + strconv.Itoa(int(frame.GetId())) + ":")
	startKey := prefix
	for {
		urlKeys, nextKey, err := db.GetPrefixP(prefix, startKey, 1000)
		if err != nil {
			fmt.Println("Error fetching subscriptions:", err)
			return
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
}

/*func Send(keys map[string][][]byte) error {
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
	} else {
		return fmt.Errorf("Error unmarshalling response body: %w", err)
	}
	return nil
}
*/
