package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Sender func(fn string, duration int, t time.Time)

// Grabber defines a procedure to grab live video
type Grabber struct {
	Segment    int    // Record time for single video
	Resolution string // 640x480... etc
	SPF        int    // seconds per frame
	Format     string
	Device     string
	Senders    []Sender
}

// Grab image from web cam
func (g *Grabber) Grab(dir string) (err error) {
	t := time.Now()
	seg := strconv.Itoa(g.Segment)
	fn := fmt.Sprintf("%s/%%0%.0f.png", dir, math.Log10(float64(g.Segment)/float64(g.SPF)))
	frameRate := fmt.Sprintf("1/%d", g.SPF)

	proc := exec.Command(
		"ffmpeg",
		"-t", seg,
		"-framerate", frameRate,
		"-f", g.Format,
		"-s", g.Resolution,
		"-i", g.Device,
		fn,
	)

	if err = proc.Run(); err != nil {
		return
	}

	go g.Encode(dir, fn, t)
	return
}

// Encode grabbed image to video
func (g *Grabber) Encode(dir, fn string, t time.Time) {
	defer g.Cleanup(dir)

	// encode to mp4. ref: http://rodrigopolo.com/ffmpeg/cheats.php
	proc := exec.Command(
		"ffmpeg",
		"-i", fn,
		"-s", g.Resolution,
		"-r", "24000/1001",
		"-b", "200k",
		"-bt", "240k",
		"-c:v", "libx264",
		"-coder", "0",
		"-bf", "0",
		"-refs", "1",
		"-flags2",
		"-wpred-dct8x8",
		"-level", "13",
		"-maxrate", "10M",
		"-bufsize", "10M",
		"-c:a", "libfaac",
		"-ac", "2",
		"-ar", "48000",
		"-ab", "192k",
		dir+"/ant.mp4",
	)

	if err := proc.Run(); err != nil {
		return
	}

	totalFrames := g.Segment / g.SPF
	duration := float64(totalFrames) / (24000.0 / 1001.0)

	g.Send(dir, int(duration), t)
}

// Cleanup work temp
func (g *Grabber) Cleanup(dir string) {
	os.RemoveAll(dir)
}

// Send video via registered senders
func (g *Grabber) Send(dir string, duration int, t time.Time) {
	for _, sender := range g.Senders {
		sender(dir + "/ant.mp4", duration, t)
	}
}
