package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	debug         = false
	defaultToken  = ""
	defaultChatID = ""
	parseMode     = "MarkdownV2" // Can be "Markdown", "MarkdownV2", "HTML", or "" for plain text
)

type Config struct {
	Token  string
	ChatID string
}

// SetDefaultToken sets the bot token to use for sending messages
func SetDefaultToken(token string) {
	defaultToken = token
}

// SetDefaultChatID sets the chat ID to send messages to
func SetDefaultChatID(chatID string) {
	defaultChatID = chatID
}

// SetDebugMode toggles verbose logs
func SetDebugMode(value bool) {
	debug = value
}

// Send sends a message to Telegram using the bot API.
// If a Config is provided, its token and chatID will be used.
// Otherwise, the default values set with SetDefaultToken and SetDefaultChatID will be used.
//
// Usage:
//
//	telegram.Send("Hello, world!") // Uses default bot token and chat ID set with SetDefaultToken and SetDefaultChatID
//	telegram.Send("Hello, world!", telegram.Config{Token: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz", ChatID: "-1001234567890"}) // Uses specified config
func Send(message string, config ...Config) {
	var botToken, chatID string
	// Use provided config values if available, otherwise use defaults
	if len(config) > 0 {
		if config[0].Token != "" {
			botToken = config[0].Token
		} else {
			botToken = defaultToken
		}
		if config[0].ChatID != "" {
			chatID = config[0].ChatID
		} else {
			chatID = defaultChatID
		}
	} else {
		botToken = defaultToken
		chatID = defaultChatID
	}
	if len(botToken) == 0 || len(chatID) == 0 {
		if debug {
			log.Println("[lazywebhooks.telegram] Bot token or chat ID not specified, skipping")
		}
		return
	}
	targetURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data, err := json.Marshal(payload{
		ChatID:    chatID,
		Text:      message,
		ParseMode: parseMode,
	})
	if err != nil {
		log.Println("[lazywebhooks.telegram] Error marshalling payload:", err.Error())
		return
	}
	request, _ := http.NewRequest("POST", targetURL, bytes.NewReader(data))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println("[lazywebhooks.telegram] Error executing request:", err.Error())
		return
	}
	defer response.Body.Close()
	if debug {
		if response.StatusCode >= 300 {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("[lazywebhooks.telegram] Error reading response body:", err.Error())
				return
			}
			log.Println("[lazywebhooks.telegram] Non-2xx response code:", string(body))
		}
	}
}

type payload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}
