package discordx

import (
	"log"
	"os"
	"strings"
)

func SendFromEnv(envKey, content string) error {
	webhookURL := strings.TrimSpace(os.Getenv(envKey))
	if webhookURL == "" {
		return nil
	}
	return New(webhookURL).Send(content)
}

func SendFromEnvAsync(envKey, content string) {
	go func() {
		if err := SendFromEnv(envKey, content); err != nil {
			log.Printf("discord send failed for %s: %v", envKey, err)
		}
	}()
}
