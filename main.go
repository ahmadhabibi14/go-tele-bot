package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// Store data from the Api
type Response []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic,omitempty"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string        `json:"definition"`
			Synonyms   []string      `json:"synonyms"`
			Antonyms   []interface{} `json:"antonyms"`
			Example    string        `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string      `json:"synonyms"`
		Antonyms []interface{} `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

func sendMessage(chatID int64, text string) error {

	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot+YOUR_TOKEN/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil

}

// This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	chat_id := body.Message.Chat.ID
	text := body.Message.Text

	log.Info("body handle: \n")
	log.Info(body)
	log.Info("TEXT: " + text)
	if strings.Contains(strings.ToLower(text), "author") {
		err := author(chat_id)
		if err != nil {
			fmt.Println("error in sending reply:", err)
		}
	} else if text != "" {

		wordsPlay(chat_id, text)
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

//The below code deals with the process of sending a response message
// to the user

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func author(chatID int64) error {

	err := sendMessage(chatID, "Alberto Vilas ; Lisbon - Portugal")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func getWordData(url string) (Response, error) {

	var result Response

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	body, err := ioutil.ReadAll(res.Body) // response body is []byte
	if err != nil {
		return result, err
	}

	// read json data into a Result struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func formatResponse(response Response) string {
	var str_response string
	var elem_str string

	for _, element := range response[0].Meanings {
		elem_str = element.PartOfSpeech + " : " + element.Definitions[0].Definition + "\n"
		str_response += elem_str

	}
	log.Info("RESPONSE: ")
	log.Info(str_response)
	return str_response
}

func wordsPlay(chatID int64, string_input string) error {
	requestURL := "https://api.dictionaryapi.dev/api/v2/entries/en/" + string_input
	response, err := getWordData(requestURL)
	if err != nil {
		sendMessage(chatID, "Meaning not found for the word: "+string_input)
		return err
	}

	// Return the 1st definition that comes for each part of speech

	err = sendMessage(chatID, formatResponse(response))
	return err
}

func log_init() error {
	file, err := os.OpenFile("log_file.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	return nil
}

func main() {

	//init the log file
	err := log_init()
	if err != nil {
		os.Exit(1)
	}

	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
