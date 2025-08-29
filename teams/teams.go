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
// Note: Teams workflows don't directly support custom usernames like Discord/Slack,
// but this is kept for API consistency. The username can be used in workflow logic if needed.
func SetUsername(user string) {
	username = user
}

// SetDefaultWebhookURL sets the Power Automate workflow URL to send messages to.
// This should be a "When a HTTP request is received" trigger URL from Power Automate,
// which replaces the deprecated Office 365 Connector webhooks.
func SetDefaultWebhookURL(url string) {
	defaultWebhookURL = url
}

// SetDebugMode toggles verbose logs
func SetDebugMode(value bool) {
	debug = value
}

// Send sends a message to a Teams workflow via Power Automate HTTP trigger.
// Uses webhookURL if specified, or defaultWebhookURL if not.
//
// The webhook URL should be from a Power Automate "When a HTTP request is received" trigger.
// This is the modern replacement for deprecated Office 365 Connector webhooks.
//
// Usage:
//
//	teams.Send("Hello, world!") // Uses default workflow URL set with SetDefaultWebhookURL
//	teams.Send("Hello, world!", "https://prod-xx.westus.logic.azure.com:443/workflows/.../triggers/manual/paths/invoke?...")
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
	
	// Create payload with message and optional username for workflow processing
	payloadData := payload{
		Text: message,
	}
	if username != "" {
		payloadData.Username = username
	}
	
	data, err := json.Marshal(payloadData)
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
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
}