package discord

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
// If none are specified, the username specified by the webhook will be used.
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
//	discord.Send("Hello, world!") // Assuming SetDefaultWebhookURL has been called first, this will use the default webhook URL (otherwise nothing will happen)
//	discord.Send("Hello, world!", "https://discord.com/api/webhooks/123456789/A1b2C3d4E5f6G7h8I9k0L1m2N3o4P5q6R7s8T9u0V1w2X3y4Z5a6B7c8D9e0F1g2")
func Send(message string, webhookURL ...string) {
	var targetURL string
	if len(webhookURL) > 0 {
		targetURL = webhookURL[0]
	} else {
		targetURL = defaultWebhookURL
	}
	if len(targetURL) == 0 {
		if debug {
			log.Println("[lazywebhooks.discord] No webhook URL specified, skipping")
		}
		return
	}
	data, err := json.Marshal(payload{Username: username, Content: message})
	if err != nil {
		log.Println("[lazywebhooks.discord] Error marshalling payload:", err.Error())
		return
	}
	request, _ := http.NewRequest("POST", targetURL, bytes.NewReader(data))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println("[lazywebhooks.discord] Error executing request:", err.Error())
		return
	}
	defer response.Body.Close()
	if debug {
		if response.StatusCode >= 300 {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("[lazywebhooks.discord] Error reading response body:", err.Error())
				return
			}
			log.Println("[lazywebhooks.discord] Non-2xx response code:", string(body))
		}
	}
}

type payload struct {
	Username string `json:"username,omitempty"`
	Content  string `json:"content"`
}
