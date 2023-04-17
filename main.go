package main

import (
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"os"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	

)


func Get(url string, queryParams string) (string) {

	msg := ""

	resp, err := http.Get(url+queryParams)

	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return "Err"
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "Err"
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return "Err"
	}

	
	var builder strings.Builder
	for key, value := range data {
		builder.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	concatString := builder.String()
	r := msg + concatString
	return r
}


func main() {

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config")
	}

	tgAPI := os.Getenv("TGAPI")
	URL := os.Getenv("URL")

	var baseURLs map[int64]string
	baseURLs = make(map[int64]string)
	
	bot, err := tgbotapi.NewBotAPI(tgAPI)
	if err != nil {
		log.Panic(err)
	}
	

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhook(URL+bot.Token)
	if err != nil {
		log.Fatal(err)
	}


	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telergam callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:443",nil)

	for update := range updates {
		
		if update.Message == nil {
			return
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			userID := update.Message.From.ID
			switch update.Message.Command() {
			case "setURL":
				baseURL := update.Message.CommandArguments()
				baseURLs[userID] = baseURL
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf("Base URL set to: %s", baseURL))
				bot.Send(msg)
			case "getURL":
				if baseURL, ok := baseURLs[userID]; ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						fmt.Sprintf("Your current URL is set to: %s", baseURL))
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No URL set")
					bot.Send(msg)
				}
			case "get":
				if baseURL, ok := baseURLs[userID]; ok {
					queryString := update.Message.CommandArguments()
					txt := Get(baseURL,queryString)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No URL set")
					bot.Send(msg)
				}
			}
		}
	}
}
