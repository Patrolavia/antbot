package main

import (
	"log"
	"os/exec"
	"time"
)

type ScpSender struct {
	Path string
	*log.Logger
}

func (s *ScpSender) Send(fn string, duration int, t time.Time) {
	proc := exec.Command("scp", fn, s.Path)

	output, err := proc.CombinedOutput()
	if err != nil {
		s.Printf("Error sending via scp to %s: %s", s.Path, err)
		s.Print(string(output))
		return
	}

	s.Printf("Successfully sent to %s via scp", s.Path)
}
