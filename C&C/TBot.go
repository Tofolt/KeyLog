package main

import (
	"github.com/buger/jsonparser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var bot *tgbotapi.BotAPI
var chatID int64

func ChatIdFromCache() int64 {
	var res int64
	jsonFile, _ := os.Open("config.json")
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	res, _ = jsonparser.GetInt(jsonBytes, "telegram_chat_ID")
	return res
}

func BotInit(token string) {
	bot, _ = tgbotapi.NewBotAPI(token)

	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)

	if ChatIdFromCache() != 0 {
		chatID = ChatIdFromCache()
		log.Println("CHAT_ID ALREADY EXISTS")
		log.Println("LISTENING FOR DATA")
	}

	for chatID == 0 {
		time.Sleep(1000 * time.Millisecond)
		for update := range updates {
			if update.Message != nil { // If we got a message
				if update.Message.Text == "/start" {
					chatID = update.Message.Chat.ID
					log.Println("GOT CHAT_ID")
					log.Println("LISTENING FOR DATA")
					MakeConfigFile()
					break
				}
			}
		}
	}
}

func TelegramBotSend(path string) {
	doc := tgbotapi.FilePath(path)
	msg := tgbotapi.NewDocument(chatID, doc)
	bot.Send(msg)
}
