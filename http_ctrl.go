package main

import (
	"log"
	"net/http"
)

type HTTPController struct {
	Bind    string
	Grabber *Grabber
	Encoder *Encoder
	*log.Logger
}

func (c *HTTPController) Control(g *Grabber, e *Encoder, l *log.Logger) {
	c.Logger = l
	c.Grabber = g
	c.Encoder = e

	http.HandleFunc("/quit", c.Quit)
	http.HandleFunc("/forcequit", c.ForceQuit)
	http.ListenAndServe(c.Bind, nil)
}

func (c *HTTPController) ForceQuit(w http.ResponseWriter, r *http.Request) {
	c.Print("[WEBC] Got force quit command")
	c.Grabber.Interrupt()
	c.Encoder.Interrupt()
}

func (c *HTTPController) Quit(w http.ResponseWriter, r *http.Request) {
	c.Print("[WEBC] Got quit command")
	c.Grabber.Interrupt()
	close(c.Encoder.Queue)
}
