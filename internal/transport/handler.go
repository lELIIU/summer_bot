package transport

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"summer_bot/internal/client"
	"summer_bot/internal/config"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{API: bot}, nil
}

func (b *Bot) Start() {
	slog.Info("bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		slog.Info("Пользователь выполнил команду",
			"userID", update.Message.From.ID,
			"userName", update.Message.From.UserName,
			"text", update.Message.Text,
			"command", update.Message.Command(),
		)

		switch update.Message.Command() {
		case "start":
			b.HandlerCommandStart(update.Message)
		case "help":
			b.HandlerCommandHelp(update.Message)
		case "weather":
			b.HandlerCommandWeather(update.Message)
		default:
			b.HandlerOtherCommand(update.Message)
		}
	}
}

func (b *Bot) HandlerCommandStart(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Для получения списка команд используй /help")
	if _, err := b.API.Send(reply); err != nil {
		slog.Error("error sending message:", err)
	}
}

func (b *Bot) HandlerCommandHelp(msg *tgbotapi.Message) {
	text := "Список команд:\n/start - запуск бота\n/help - получить список команд\n/weather - получение прогноза погоды(необходимо ввести название локации на английском)"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	if _, err := b.API.Send(reply); err != nil {
		slog.Error("error sending message:", err)
	}
}

func (b *Bot) HandlerCommandWeather(msg *tgbotapi.Message) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config in handler.")
	}

	weather := client.NewWeatherClient(cfg.WeatherApiKey)

	args := strings.Fields(msg.Text)

	if len(args) < 2 {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "Используйте: /weather <город> [дней=1]")
		b.API.Send(resp)
		return
	}

	city := args[1] // Лондон
	days := 1       // Значение по умолчанию

	// Если пользователь указал количество дней
	if len(args) >= 3 {
		if d, err := strconv.Atoi(args[2]); err == nil {
			days = d
		}
	}
	if msg.CommandArguments() == "" {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Используйте: /weather <город> [дней=1]")
		if _, err := b.API.Send(reply); err != nil {
			slog.Error("error sending message:", err)
		}
	}

	resp, err := weather.GetCurrentWeather(city, days)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, resp)

	if _, err := b.API.Send(reply); err != nil {
		slog.Error("error sending message:", err)
	}
}

func (b *Bot) HandlerOtherCommand(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "на данный момент бот принимает только 1 команду. Посмотреть их можно прописав /help")
	if _, err := b.API.Send(reply); err != nil {
		slog.Error("error sending message:", err)
	}
}
