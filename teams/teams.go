package teams

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	debug             = false
	defaultWebhookURL = ""
	username          = ""
)

// SetUsername sets the username the message will be sent with.
// Note: Microsoft Teams webhooks don't support custom usernames in the same way as Discord/Slack,
// but this is kept for API consistency across all webhook providers.
func SetUsername(user string) {
	username = user
}

// SetDefaultWebhookURL sets the webhook URL to send messages to
func SetDefaultWebhookURL(url string) {
	defaultWebhookURL = url
}

// SetDebugMode toggles verbose logs
func SetDebugMode(value bool) {
	debug = value
}

// Send sends a message to webhookURL if it is specified, or defaultWebhookURL if it is not.
//
// Usage:
//
//	teams.Send("Hello, world!") // Assuming SetDefaultWebhookURL has been called first, this will use the default webhook URL (otherwise nothing will happen)
//	teams.Send("Hello, world!", "https://outlook.office.com/webhook/...")
func Send(message string, webhookURL ...string) {
	var targetURL string
	if len(webhookURL) > 0 {
		targetURL = webhookURL[0]
	} else {
		targetURL = defaultWebhookURL
	}
	if len(targetURL) == 0 {
		if debug {
			log.Println("[lazywebhooks.teams] No webhook URL specified, skipping")
		}
		return
	}
	data, err := json.Marshal(payload{Text: message})
	if err != nil {
		log.Println("[lazywebhooks.teams] Error marshalling payload:", err.Error())
		return
	}
	request, _ := http.NewRequest("POST", targetURL, bytes.NewReader(data))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println("[lazywebhooks.teams] Error executing request:", err.Error())
		return
	}
	defer response.Body.Close()
	if debug {
		if response.StatusCode >= 300 {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("[lazywebhooks.teams] Error reading response body:", err.Error())
				return
			}
			log.Println("[lazywebhooks.teams] Non-2xx response code:", string(body))
		}
	}
}

type payload struct {
	Text string `json:"text"`
}