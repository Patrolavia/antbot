package main

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type CLIController struct {
	Input io.Reader
}

func (c *CLIController) Control(g *Grabber, e *Encoder, l *log.Logger) {
	r := bufio.NewReader(c.Input)
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			// unable to use cli controller, exit.
			l.Print("[CLIC] Unable read data from stdin, cli controller disabled.")
			return
		}
		l.Printf("[CLIC] Got input %s", str)
		switch strings.ToLower(strings.TrimSpace(str)) {
		case "q", "quit":
			l.Print("[CLIC] Got quit command")
			g.Interrupt()
			return
		case "fq", "forcequit":
			l.Print("[CLIC] Got force quit command")
			g.Interrupt()
			e.Interrupt()
		}
	}
}
