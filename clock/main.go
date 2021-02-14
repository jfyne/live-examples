package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/jfyne/live"
)

const (
	tick = "tick"
)

type clock struct {
	Time time.Time
}

func newClock(s *live.Socket) *clock {
	c, ok := s.Assigns().(*clock)
	if !ok {
		return &clock{
			Time: time.Now(),
		}
	}
	return c
}

func (c clock) FormattedTime() string {
	return c.Time.Format("15:04:05")
}

func mount(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
	// Take the socket data and tranform it into our view model if it is
	// available.
	c := newClock(s)

	// If we are mouting the websocket connection, trigger the first tick
	// event.
	if s.Connected() {
		go func() {
			time.Sleep(1 * time.Second)
			s.Self(ctx, live.Event{T: tick})
		}()
	}
	return c, nil
}

func main() {
	t, err := template.ParseFiles("root.html", "clock/view.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	// Set the mount function for this handler.
	h.Mount = mount

	// Server side events.

	// tick event updates the clock every second.
	h.HandleSelf(tick, func(ctx context.Context, s *live.Socket, _ map[string]interface{}) (interface{}, error) {
		// Get our model
		c := newClock(s)
		// Update the time.
		c.Time = time.Now()
		// Send ourselves another tick in a second.
		go func(sock *live.Socket) {
			time.Sleep(1 * time.Second)
			s.Self(ctx, live.Event{T: tick})
		}(s)
		return c, nil
	})

	// Run the server.
	http.Handle("/clock", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
