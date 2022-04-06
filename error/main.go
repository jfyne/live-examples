package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jfyne/live"
)

const (
	problem = "problem"
)

func main() {
	t, err := template.ParseFiles("root.html", "error/view.html")
	if err != nil {
		log.Fatal(err)
	}

	h := live.NewHandler(live.WithTemplateRenderer(t))

	// Uncomment the below to see the server respond with an error immediately.

	//h.HandleMount(func(ctx context.Context, s live.Socket) (interface{}, error) {
	//	return nil, fmt.Errorf("mount failure")
	//})

	h.HandleError(func(ctx context.Context, err error) {
		w := live.Writer(ctx)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("this is a bad request: " + err.Error()))
	})

	h.HandleEvent(problem, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return nil, fmt.Errorf("hello")
	})

	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
