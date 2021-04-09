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

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	h.Error = func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("this is a bad request"))
	}

	h.HandleEvent(problem, func(ctx context.Context, s *live.Socket, _ live.Params) (interface{}, error) {
		return nil, fmt.Errorf("hello")
	})

	http.Handle("/error", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
