package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Patrolavia/botgoram/telegram"
)

func main() {
	var (
		tokenFile  string
		channel    string
		segment    int
		resolution string
		spf        int
		format     string
		device     string
		bot        telegram.API
	)

	flag.StringVar(&tokenFile, "t", "token", "The file holding telegram bot token")
	flag.StringVar(&channel, "ch", "", "Telegram channel to announce your video, leave empty if not using this feature")
	flag.IntVar(&segment, "seg", 1800, "Time to record for each video segment")
	flag.StringVar(&resolution, "size", "640x480", "Cam source resolution")
	flag.IntVar(&spf, "spf", 1, "Grab 1 frame every `N` seconds")
	flag.StringVar(&format, "f", "v4l2", "Source format for ffmpeg")
	flag.StringVar(&device, "i", "/dev/video0", "Input file for ffmpeg")
	flag.Parse()

	logFile, err := os.Create("ant.log")
	if err != nil {
		log.Fatalf("Cannot create logfile ant.log: %s", err)
	}

	logger := log.New(io.MultiWriter(logFile, os.Stderr), "", log.LstdFlags)

	senders := []Sender{}

	if channel != "" {
		bot = initTelegram(tokenFile)
		senders = append(senders, &TelegramChannelSender{bot, channel, logger})
	}

	i := 0

	grabber := &Grabber{
		Segment:    segment,
		Resolution: resolution,
		SPF:        spf,
		Format:     format,
		Device:     device,
		Senders:    senders,
		Logger:     logger,
	}

	for {
		dir := "work" + strconv.Itoa(i)
		if err := os.Mkdir(dir, 0644); err != nil {
			logger.Fatalf("Cannot create temp dir %s: %s", dir, err)
		}
		grabber.Grab(dir)
		if i++; i > 99 {
			i = 0
		}
	}
}
