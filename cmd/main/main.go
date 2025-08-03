package main

import (
	"fmt"
	"log"
	"summer_bot/internal/config"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config.")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "вы написали "+update.Message.Chat.UserName)

		if _, err := bot.Send(msg); err != nil {
			log.Println("Ошибка отправки:", err)
		}

	}
}
