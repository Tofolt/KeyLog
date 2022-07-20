package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatID int64

func BotInit() {
	bot, _ = tgbotapi.NewBotAPI("<Token>")

	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == "/start" {
				chatID = update.Message.Chat.ID
			}
		}
	}
}

func TelegramBotSend(path string) {
	bot.Debug = true

	doc := tgbotapi.FilePath(path)
	msg := tgbotapi.NewDocument(chatID, doc)
	bot.Send(msg)
}
