package main

import (
	"context"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/components"
	"github.com/jfyne/live/page"
)

func main() {
	// Setup handler.
	h := live.NewHandler(
		page.WithComponentMount(func(ctx context.Context, h live.Handler, s live.Socket) (*page.Component, error) {
			return components.NewPage("app", h, s, "Clocks")
		}),
		page.WithComponentRenderer(),
	)

	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
