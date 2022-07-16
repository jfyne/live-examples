package main

import (
	"context"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/components"
	"github.com/jfyne/live/page"
)

func main() {
	h := page.NewHandler(func(ctx context.Context, h *live.Handler, s live.Socket) (page.ComponentLifecycle, error) {
		return components.NewClocks("Clocks")
	})

	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
