# lazywebhooks
Webhooks for lazy Go devs.


## discord
```go
package main

import (
	"github.com/TwiN/lazywebhooks/discord"
)

func main() {
	// You can set a username if you wish. If you don't, it will default to your webhook's username
	discord.SetUsername("SomeUsername")
	// You can set a default webhook URL if you don't want to include the webhook URL in every discord.Send(...) call you make
	discord.SetDefaultWebhookURL("https://discord.com/api/webhooks/1234567890/A1b2C3d4E5f6G7h8I9k0L1m2N3o4P5q6R7s8T9u0V1w2X3y4Z5a6B7c8D9e0F1g2")
	// If you don't specify a webhook URL in the discord.Send(...) call, it will use the default webhook URL you set above
	discord.Send("Hello, world!")
	discord.Send("My name is John Doe")
	// Otherwise, if you do set the webhookURL parameter in the discord.Send(...) call, it will use that instead of the default webhook URL
	discord.Send("Hello, world!", "https://discord.com/api/webhooks/0987654321/z9Y8x7W6v5U4t3S2r1Q0p9O8n7M6l5K4j3I2h1G0f9D8e7C6b5A4z3X2c1V0b9N8")
}
```
