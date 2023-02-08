package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func main() {
	// Create a new OpenAI API client
	c := gogpt.NewClient(os.Getenv("OPENAI_API_TOKEN"))
	ctx := context.Background()

	// Create a new Telegram bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %s", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start listening for updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30000

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Use the OpenAI API client to generate a response
		req := gogpt.CompletionRequest{
			Model:     gogpt.GPT3Davinci, /*gogpt.GPT3Ada*/
			MaxTokens: 5,
			Prompt:    update.Message.Text,
		}

		resp, err := c.CreateCompletion(ctx, req)
		if err != nil {
			log.Printf("Failed to generate response from OpenAI: %s", err)
			continue
		}

		// Format the response as a string
		// responseText := fmt.Sprintf("Response from OpenAI's GPT-3: %s", resp)

		// Send the response back to the user
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send response to user: %s", err)
			continue
		}
	}
}
