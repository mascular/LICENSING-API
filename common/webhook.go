package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

func SendUsageAlert(app, key, duration, action string) {
	if AppConfig.DiscordWebhook == "" {
		fmt.Println("⚠️ No Discord webhook configured")
		return
	}

	payload := map[string]string{
		"content": "**Usage Alert | License**\n" +
			"> App: `" + app + "`\n" +
			"> Key: `" + key + "`\n" +
			"> Duration: `" + duration + "`\n" +
			"> Action: `" + action + "`",
	}

	body, _ := json.Marshal(payload)
	http.Post(AppConfig.DiscordWebhook, "application/json", bytes.NewBuffer(body))
}
