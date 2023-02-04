package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30000

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Ignore any non-message updates
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // Ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet.
		// so we leave it empty
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know the command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
