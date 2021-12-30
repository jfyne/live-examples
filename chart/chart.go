package main

import (
	"context"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jfyne/live"
)

const (
	regenerate = "regenerate"
)

type chartData struct {
	Sales []int
}

func newChartData(s live.Socket) *chartData {
	d, ok := s.Assigns().(*chartData)
	if !ok {
		return &chartData{
			Sales: rand.Perm(9),
		}
	}
	return d
}

func main() {
	rand.Seed(time.Now().Unix())

	t, err := template.ParseFiles("root.html", "chart/view.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	// Set the mount function for this handler.
	h.HandleMount(func(ctx context.Context, s live.Socket) (interface{}, error) {
		// This will initialise the chart data if needed.
		return newChartData(s), nil
	})

	// Client side events.

	// Regenerate event, creates new random sales data.
	h.HandleEvent(regenerate, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		// Get this sockets counter struct.
		c := newChartData(s)

		// Generate new sales data.
		c.Sales = rand.Perm(9)

		// Set the new chart data back to the socket.
		return c, nil
	})

	// Run the server.
	http.Handle("/chart", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
