package main

import (
	"fmt"
	"log"
	"net/http"
)

func html(title, body string) []byte {
	return []byte(fmt.Sprintf(
		`<html><head><title>%s</title><meta charset="utf8">
<style>.controller {position: absolute; display: block; height: 200px; width: 200px; top: 50%%; left: 50%%; margin: -150px 0 0 -100px; } .button {display: block; height: 80px; width: 180px; line-height: 80px; border-redius: 10px; margin: 10px; font-size: 42px; text-shadow: 0 -1px 1px rgba(255, 255, 255, .5); } .button-quit { color: #E53935; } .button-forceQuit { color: #D81B60; }</style>
</head><body>%s</body></html>`,
		title,
		body,
	))
}

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

	http.HandleFunc("/quit", c.quit)
	http.HandleFunc("/forcequit", c.forceQuit)
	http.HandleFunc("/", c.index)
	http.ListenAndServe(c.Bind, nil)
}

func (c *HTTPController) forceQuit(w http.ResponseWriter, r *http.Request) {
	c.Print("[WEBC] Got force quit command")
	c.Grabber.Interrupt()
	c.Encoder.Interrupt()
	w.Write(html(
		`Shutdown antbot`,
		`Stop antbot now.`,
	))
}

func (c *HTTPController) quit(w http.ResponseWriter, r *http.Request) {
	c.Print("[WEBC] Got quit command")
	c.Grabber.Interrupt()
	w.Write(html(
		`Shutdown antbot`,
		`Encoding last video, bot will go down ater sending last video.`,
	))
}

func (c *HTTPController) index(w http.ResponseWriter, r *http.Request) {
	w.Write(html(
		`Antbot web control`,
		`<div class="controller"><a href="/quit" class="button button-quit">QUIT</a><a href="/forcequit" class="button button-forceQuit">Force quit</a><div>`,
	))
}
