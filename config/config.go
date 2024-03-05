package config

import (
	"os"

	"github.com/caarlos0/env"
	gotdotenv "github.com/joho/godotenv"
)

func LoadConfigFromEnv(c any) error {
	if os.Getenv("ENV") != "production" {
		if err := gotdotenv.Load(); err != nil {
			return err
		}
	}

	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}
