package main

import (
	"fmt"
	"log"
	"net/http"
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
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Retrieve response from OpenAI's GPT-3 language model
		response, err := http.Get("https://api.openai.com/v1/engines/gpt-3/jobs")
		if err != nil {
			log.Panic(err)
		}

		// Format response as a string
		responseText := fmt.Sprintf("Response from GPT-3: %s", response)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
