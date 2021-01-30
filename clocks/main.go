package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/components"
	"github.com/jfyne/live/component"
)

func main() {
	// Setup handler.
	h, err := live.NewHandler(
		live.NewCookieStore("session-name", []byte("weak-secret")),
		component.WithComponentMount(func(ctx context.Context, h *live.Handler, r *http.Request, s *live.Socket) (component.Component, error) {
			return components.NewPage("app", h, s, "Example page")
		}),
		component.WithComponentRenderer(),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/clocks", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
