package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type Sender interface {
	Send(fn string, duration int, t time.Time)
}

type WorkInfo struct {
	Dir      string
	Filename string
	Time     time.Time
}

// Grabber defines a procedure to grab live video
type Grabber struct {
	Segment    int    // Record time for single video
	Resolution string // 640x480... etc
	SPF        int    // seconds per frame
	Format     string
	Device     string
	process    *os.Process
	*sync.Mutex
	*log.Logger
}

func (g *Grabber) Interrupt() {
	proc := g.process

	if proc != nil {
		proc.Kill()
	}
}

// Grab image from web cam
func (g *Grabber) Grab(dir string) (ret WorkInfo, err error) {
	g.Lock()
	defer g.Unlock()
	t := time.Now()
	seg := strconv.Itoa(g.Segment)
	fn := fmt.Sprintf("%s/%%0%.0fd.png", dir, math.Ceil(math.Log10(float64(g.Segment)/float64(g.SPF))))
	g.Printf("[GRAB] Grabbing to dir %s ...", fn)
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

	g.process = proc.Process
	err = proc.Run()
	g.process = nil

	ret = WorkInfo{
		Dir:      dir,
		Filename: fn,
		Time:     t,
	}

	return
}

type Encoder struct {
	Segment    int    // Record time for single video
	Resolution string // 640x480... etc
	SPF        int    // seconds per frame
	Senders    []Sender
	Queue      chan WorkInfo
	process    *os.Process
	*log.Logger
}

func (e *Encoder) Interrupt() {
	close(e.Queue)
	for range e.Queue {
	}
	proc := e.process
	if proc != nil {
		proc.Kill()
	}
}

// Encode grabbed image to video
func (e *Encoder) Run() {
	for work := range e.Queue {
		fn, dir, t := work.Filename, work.Dir, work.Time
		if fn == "" || dir == "" {
			e.Print("[ENC ] Got empty workinfo")
			continue
		}

		if _, err := os.Stat(fmt.Sprintf(fn, 1)); err != nil {
			e.Printf("[ENC ] No grabbed data found in %s ...", dir)
			continue
		}
		e.Printf("[ENC ] Encoding from %s ...", fn)

		// encode to mp4. ref: http://rodrigopolo.com/ffmpeg/cheats.php
		proc := exec.Command(
			"ffmpeg",
			"-i", fn,
			"-s", e.Resolution,
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

		f, err := os.Create(dir + "/encode.out")
		if err != nil {
			e.Printf("[ENC ] Cannot create log file for subprocess: %s", err)
			continue
		}
		proc.Stdout = f
		proc.Stderr = f
		defer f.Close()

		e.process = proc.Process
		if err = proc.Run(); err != nil {
			e.Printf("[ENC ] Subprocess returns error: %s", err)
			continue
		}
		e.process = nil

		totalFrames := e.Segment / e.SPF
		duration := float64(totalFrames) / (24000.0 / 1001.0)

		e.Send(dir, int(duration), t)
		e.Cleanup(work)
	}
	e.Print("[ENC ] No more works.")
}

// Cleanup work temp
func (e *Encoder) Cleanup(work WorkInfo) {
	e.Printf("[ENC ] Cleaning up %s ...", work.Dir)
	os.RemoveAll(work.Dir)
}

// Send video via registered senders
func (e *Encoder) Send(dir string, duration int, t time.Time) {
	for _, sender := range e.Senders {
		sender.Send(dir+"/ant.mp4", duration, t)
	}
}
