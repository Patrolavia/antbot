package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Patrolavia/botgoram/telegram"
)

func main() {
	tokenByte, err := ioutil.ReadFile("token")
	if err != nil {
		log.Fatalf(`Cannot read telegram bot token from file "token": %s`, err)
	}

	token := strings.TrimSpace(string(tokenByte))

	bot := telegram.New(token)

	if _, err := bot.Me(); err != nil {
		log.Fatalf("Error validating bot: %s", err)
	}

	i := 0

	grabber := &Grabber{
		Segment:    1800,
		Resolution: "640x480",
		SPF:        1,
		Format:     "v4l2",
		Device:     "/dev/video0",
		Senders: []Sender{func(fn string, duration int, t time.Time) {
			caption := t.Local().Format("2006-1-2 15:04:05")
			bot.SendVideo(
				&telegram.Chat{User: &telegram.User{Username: "@ronmiants"}, Type: "channel"},
				&telegram.File{Filename: fn, MimeType: "video/mp4"},
				duration,
				caption,
				nil,
			)
		}},
	}

	for {
		dir := strconv.Itoa(i)
		if err := os.Mkdir(dir, 0644); err != nil {
			log.Fatalf("Cannot create temp dir %s: %s", dir, err)
		}
		grabber.Grab(dir)
	}
}
