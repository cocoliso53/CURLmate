package main

import (
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

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
	bot, err := tgbotapi.NewBotAPI("")

	var baseURLs map[int64]string

	baseURLs = make(map[int64]string)
	

	if err != nil {
		log.Panic(err)
	}


	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

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
