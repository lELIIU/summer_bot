package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TelegramBotToken string
	WeatherApiKey    string
	CryptApiKey      string
	StockApiKey      string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found, relying on environment variables")
	}
	return &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_API_TOKEN"),
		WeatherApiKey:    os.Getenv("WEATHER_API_TOKEN"),
		CryptApiKey:      os.Getenv("CRYPT_API_TOKEN"),
		StockApiKey:      os.Getenv("STOCK_API_TOKEN"),
	}, nil
}
