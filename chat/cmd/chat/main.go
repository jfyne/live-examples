package main

import (
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/chat"
)

func main() {
	// Run the server.
	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), chat.NewHandler()))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
