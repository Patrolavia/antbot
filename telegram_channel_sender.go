package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/Patrolavia/botgoram/telegram"
)

type TelegramChannelSender struct {
	API     telegram.API
	Channel string
}

func (s *TelegramChannelSender) Send(fn string, duration int, t time.Time) {
	caption := t.Local().Format("2006-1-2 15:04:05")
	s.API.SendVideo(
		&telegram.Chat{User: &telegram.User{Username: s.Channel}, Type: telegram.TYPECHANNEL},
		&telegram.File{Filename: fn, MimeType: "video/mp4"},
		duration,
		caption,
		nil,
	)
}

func initTelegram(file string) telegram.API {
	tokenByte, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf(`Cannot read telegram bot token from file "token": %s`, err)
	}

	bot := telegram.New(strings.TrimSpace(string(tokenByte)))

	if _, err := bot.Me(); err != nil {
		log.Fatalf("Error validating bot: %s", err)
	}
	return bot
}
