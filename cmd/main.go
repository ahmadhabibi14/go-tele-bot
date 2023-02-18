package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ahmadhabibi14/go-tele-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Initialize .env file
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %s", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		reqBody := models.Request{
			ModelRequest: "text-davinci-003",
			Prompt:       update.Message.Text,
			Temperature:  1,
			MaxTokens:    100,
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonBody))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
		if err != nil {
			panic(err)
		}

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Read the response
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var data models.TextCompletionResponse
		json.Unmarshal([]byte(jsonData), &data)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, data.Choices[0].Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
