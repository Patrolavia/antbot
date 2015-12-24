package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Sender interface {
	Send(fn string, duration int, t time.Time)
}

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
func (g *Grabber) Grab(dir string) {
	t := time.Now()
	seg := strconv.Itoa(g.Segment)
	fn := fmt.Sprintf("%s/%%0%.0fd.png", dir, math.Log10(float64(g.Segment)/float64(g.SPF)))
	log.Printf("Grabbing to dir %s ...", fn)
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

	if f, err := os.Create(dir + "/grab.out"); err == nil {
		proc.Stdout = f
		proc.Stderr = f
		defer f.Close()
	}

	if err := proc.Run(); err == nil {
		go g.Encode(dir, fn, t)
	} else {
		log.Fatalf("Grabber error: %s", err)
	}
}

// Encode grabbed image to video
func (g *Grabber) Encode(dir, fn string, t time.Time) {
	defer g.Cleanup(dir)
	log.Printf("Encoding from %s ...", fn)

	// encode to mp4. ref: http://rodrigopolo.com/ffmpeg/cheats.php
	proc := exec.Command(
		"ffmpeg",
		"-i", fn,
		"-s", g.Resolution,
		"-r", "24000/1001",
		"-profile:v", "main",
		"-level", "4.0",
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-c:a", "libfaac",
		"-ac", "2",
		"-ar", "48000",
		"-ab", "192k",
		dir+"/ant.mp4",
	)

	if f, err := os.Create(dir + "/encode.out"); err == nil {
		proc.Stdout = f
		proc.Stderr = f
		defer f.Close()
	}

	if err := proc.Run(); err != nil {
		log.Fatalf("Encoder error: %s", err)
	}

	totalFrames := g.Segment / g.SPF
	duration := float64(totalFrames) / (24000.0 / 1001.0)

	g.Send(dir, int(duration), t)
}

// Cleanup work temp
func (g *Grabber) Cleanup(dir string) {
	log.Printf("Cleaning up %s ...", dir)
	os.RemoveAll(dir)
}

// Send video via registered senders
func (g *Grabber) Send(dir string, duration int, t time.Time) {
	for _, sender := range g.Senders {
		sender.Send(dir+"/ant.mp4", duration, t)
	}
}
