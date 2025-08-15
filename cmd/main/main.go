package main

import (
	"log"
	"log/slog"
	_ "log/slog"
	"os"
	"summer_bot/internal/config"
	"summer_bot/internal/transport"
	//"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config in main.")
	}

	bot, err := transport.NewBot(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Bot is authorized")

	bot.Start()

}
