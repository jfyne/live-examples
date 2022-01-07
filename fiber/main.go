package main

import (
	"context"
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/jfyne/live"
	"github.com/jfyne/live-contrib/livefiber"
)

const (
	inc = "inc"
	dec = "dec"
)

type counter struct {
	Value int
}

func newCounter(s live.Socket) *counter {
	c, ok := s.Assigns().(*counter)
	if !ok {
		return &counter{}
	}
	return c
}

func main() {
	h := live.NewHandler(livefiber.WithViewsRenderer("view", "root"))

	// Set the mount function for this handler.
	h.HandleMount(func(ctx context.Context, s live.Socket) (interface{}, error) {
		// This will initialise the counter if needed.
		return newCounter(s), nil
	})

	// Client side events.

	// Increment event. Each click will increment the count by one.
	h.HandleEvent(inc, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		// Get this sockets counter struct.
		c := newCounter(s)

		// Increment the value by one.
		c.Value += 1

		// Set the counter struct back to the socket data.
		return c, nil
	})

	// Decrement event. Each click will increment the count by one.
	h.HandleEvent(dec, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		// Get this sockets counter struct.
		c := newCounter(s)

		// Decrement the value by one.
		c.Value -= 1

		// Set the counter struct back to the socket data.
		return c, nil
	})

	// Run the server.
	app := fiber.New(fiber.Config{
		Views: html.New("./fiber", ".html"),
	})

	app.Get("/fiber", livefiber.NewHandler(session.New(), h).Handlers()...)
	app.Get("/live.js", adaptor.HTTPHandler(live.Javascript{}))
	app.Get("/auto.js.map", adaptor.HTTPHandler(live.JavascriptMap{}))

	log.Fatal(app.Listen(":8080"))

}
