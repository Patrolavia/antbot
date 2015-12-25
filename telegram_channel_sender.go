package main

import (
	"log"
	"time"

	"github.com/Patrolavia/botgoram/telegram"
)

type TelegramChannelSender struct {
	API     telegram.API
	Channel string
	*log.Logger
}

func (s *TelegramChannelSender) Send(fn string, duration int, t time.Time) {
	caption := t.Local().Format("2006-1-2 15:04:05")
	_, err := s.API.SendVideo(
		&telegram.Chat{User: &telegram.User{Username: s.Channel}, Type: telegram.TYPECHANNEL},
		&telegram.File{Filename: fn, MimeType: "video/mp4"},
		duration,
		caption,
		nil,
	)
	s.Printf("%s video sent to channel %s (err=%v)", caption, s.Channel, err)
}
