package common

import (
    "encoding/json"
    "log"
    "os"
)

type Config struct {
    APIKey string `json:"X-Api-Key"`
    Port   string `json:"Port"`
}

var AppConfig Config

func LoadConfig() {
    file, err := os.Open("config.json")
    if err != nil {
        log.Fatalf("❌ Failed to open config.json: %v", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&AppConfig)
    if err != nil {
        log.Fatalf("❌ Failed to parse config.json: %v", err)
    }

    if AppConfig.APIKey == "" || AppConfig.Port == "" {
        log.Fatalf("❌ API key and Port must not be empty in config.json")
    }

    log.Println("✅ Configuration loaded successfully")
}
