package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func TelegramBotSend(path string) {

	bot, err := tgbotapi.NewBotAPI(<Token>)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	doc := tgbotapi.FilePath(path)
	msg := tgbotapi.NewDocument(402464676, doc)
	bot.Send(msg)
}
