package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/chat"
)

func main() {
	// Run the server.
	e := live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), chat.NewHandler())
	go func() {
		for {
			e.Broadcast("newmessage", chat.Message{ID: live.NewID(), User: "Room", Msg: fmt.Sprintf("The time is now %s", time.Now().Format(time.Kitchen))})
			time.Sleep(1 * time.Minute)
		}
	}()
	http.Handle("/", e)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
