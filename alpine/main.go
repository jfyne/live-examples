package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/jfyne/live"
)

//go:embed index.js
var static embed.FS

const (
	submit   = "submit"
	suggest  = "suggest"
	selected = "selected"
)

type item struct {
	ID   string
	Name string
}

func (i item) Match(search string) bool {
	s := strings.ToLower(search)
	if strings.Contains(strings.ToLower(i.Name), s) {
		return true
	}
	if strings.Contains(strings.ToLower(i.ID), s) {
		return true
	}
	return false
}

type Autocomplete struct {
	items       []item
	Suggestions []item
	Selected    []item
}

func newAutocomplete(s live.Socket) *Autocomplete {
	a, ok := s.Assigns().(*Autocomplete)
	if !ok {
		return &Autocomplete{}
	}
	return a
}

func mount(ctx context.Context, s live.Socket) (interface{}, error) {
	a := newAutocomplete(s)
	a.items = []item{
		{ID: "1", Name: "Item One"},
		{ID: "2", Name: "Item Two"},
		{ID: "3", Name: "Item Three"},
		{ID: "4", Name: "Item Four"},
		{ID: "5", Name: "Item Five"},
	}
	return a, nil
}

func main() {
	t, err := template.ParseFiles("alpine/root.html", "alpine/view.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}

	h.HandleMount(mount)

	h.HandleEvent(suggest, func(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
		a := newAutocomplete(s)
		a.Suggestions = []item{}
		search := p.String("search")
		for _, i := range a.items {
			if i.Match(search) {
				a.Suggestions = append(a.Suggestions, i)
			}
		}
		return a, nil
	})

	h.HandleEvent(selected, func(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
		a := newAutocomplete(s)
		id := p.String("id")
		// Dont select option more than once.
		for _, i := range a.Selected {
			if i.ID == id {
				return a, nil
			}
		}
		for _, i := range a.items {
			if i.ID == id {
				a.Selected = append(a.Selected, i)
				break
			}
		}
		return a, nil
	})

	h.HandleEvent(submit, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return s.Assigns(), nil
	})

	http.Handle("/alpine", h)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	http.ListenAndServe(":8080", nil)
}
