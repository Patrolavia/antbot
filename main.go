package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

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
		logFile    string
		scpPath    string
		httpBind   string
		bot        telegram.API
	)

	flag.IntVar(&segment, "seg", 1800, "Time to record for each video segment")
	flag.StringVar(&resolution, "size", "640x480", "Cam source resolution")
	flag.IntVar(&spf, "spf", 1, "Grab 1 frame every `N` seconds")
	flag.StringVar(&format, "f", "v4l2", "Source format for ffmpeg")
	flag.StringVar(&device, "i", "/dev/video0", "Input file for ffmpeg")
	flag.StringVar(&logFile, "l", "ant.log", "Log file")
	flag.StringVar(&tokenFile, "t", "token", "The file holding telegram bot token")
	flag.StringVar(&channel, "ch", "", "Telegram channel to announce your video, leave empty if not using this feature")
	flag.StringVar(&scpPath, "scp", "", "Send video to this path via scp")
	flag.StringVar(&httpBind, "webc", "", "Enable web controller at [ip address]:port, eg. :8000 or 192.168.1.1:1234")
	flag.Parse()

	logf, err := os.Create(logFile)
	if err != nil {
		log.Fatalf("Cannot create logfile ant.log: %s", err)
	}

	logger := log.New(io.MultiWriter(logf, os.Stderr), "", log.LstdFlags)

	senders := []Sender{}

	if channel != "" {
		bot = initTelegram(tokenFile, bot)
		senders = append(senders, &TelegramChannelSender{bot, channel, logger})
	}

	if scpPath != "" {
		senders = append(senders, &ScpSender{scpPath, logger})
	}

	grabber := &Grabber{
		Segment:    segment,
		Resolution: resolution,
		SPF:        spf,
		Format:     format,
		Device:     device,
		Logger:     logger,
		Mutex:      &sync.Mutex{},
	}

	encoder := &Encoder{
		Segment:    segment,
		Resolution: resolution,
		SPF:        spf,
		Senders:    senders,
		Queue:      make(chan WorkInfo, 1),
		Logger:     logger,
	}

	logger.Print("[MAIN] Starting flow controller ...")
	ctrl := &CLIController{os.Stdin}
	go ctrl.Control(grabber, encoder, logger)

	if httpBind != "" {
		logger.Print("[MAIN] Starting web flow controller ...")
		webc := &HTTPController{Bind: httpBind}
		go webc.Control(grabber, encoder, logger)
	}

	logger.Print("[MAIN] Starting cam grabber ...")
	go func(g *Grabber, e *Encoder, l *log.Logger) {
		i := 0
		for {
			dir := "work" + strconv.Itoa(i)
			if err := os.Mkdir(dir, 0755); err != nil {
				l.Fatalf("Cannot create temp dir %s: %s", dir, err)
			}
			ret, err := g.Grab(dir)
			if err != nil {
				close(e.Queue)
				return
			}
			e.Queue <- ret
			if i++; i > 99 {
				i = 0
			}
		}
	}(grabber, encoder, logger)

	logger.Printf("[MAIN] Starting encoder ...")
	encoder.Run()
}

func initTelegram(file string, bot telegram.API) telegram.API {
	if bot != nil {
		return bot
	}
	tokenByte, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf(`[MAIN] Cannot read telegram bot token from file "token": %s`, err)
	}

	ret := telegram.New(strings.TrimSpace(string(tokenByte)))

	if _, err := ret.Me(); err != nil {
		log.Fatalf("[MAIN] Error validating bot: %s", err)
	}
	return ret
}

type Controller interface {
	Control(g *Grabber, e *Encoder, l *log.Logger)
}
