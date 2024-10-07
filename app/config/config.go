package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	TgBotHost     string
	Token         string
	StoragePath   string
	BatchSize     int
	HelpVoicePath string
}

var loadedConfig *Config

func Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	loadedConfig = &Config{
		BatchSize:     100,
		TgBotHost:     "api.telegram.org",
		Token:         os.Getenv("TG_BOT_TOKEN"),
		StoragePath:   os.Getenv("STORAGE_PATH"),
		HelpVoicePath: os.Getenv("HELP_VOICE_PATH"),
	}

	return nil
}

func Get() *Config {
	return loadedConfig
}
