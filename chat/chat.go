package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jfyne/live"
)

const (
	send       = "send"
	newmessage = "newmessage"
)

type Message struct {
	ID   string // Unique ID per message so that we can use `live-update`.
	User string
	Msg  string
}

type ChatInstance struct {
	Messages []Message
}

func NewChatInstance(s *live.Socket) *ChatInstance {
	m, ok := s.Assigns().(*ChatInstance)
	if !ok {
		return &ChatInstance{
			Messages: []Message{
				{ID: live.NewID(), User: "Room", Msg: "Welcome to chat " + s.Session.ID},
			},
		}
	}
	return m
}

func NewHandler() *live.Handler {
	t, err := template.ParseFiles("chat/layout.html", "chat/view.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}
	// Set the mount function for this handler.
	h.Mount = func(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
		// This will initialise the chat for this socket.
		return NewChatInstance(s), nil
	}

	// Handle user sending a message.
	h.HandleEvent(send, func(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
		m := NewChatInstance(s)
		msg := live.ParamString(p, "message")
		if msg == "" {
			return m, nil
		}
		data, err := json.Marshal(Message{ID: live.NewID(), User: s.Session.ID, Msg: msg})
		if err != nil {
			return m, fmt.Errorf("failed marshalling message for broadcast: %w", err)
		}
		h.Broadcast(live.Event{T: newmessage, Data: map[string]interface{}{"message": string(data)}})
		return m, nil
	})

	// Handle the broadcasted events.
	h.HandleSelf(newmessage, func(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
		m := NewChatInstance(s)
		data, ok := p["message"]
		if !ok {
			return m, fmt.Errorf("no message key")
		}
		raw, ok := data.(string)
		if !ok {
			return m, fmt.Errorf("no message bytes")
		}
		var msg Message
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			return m, fmt.Errorf("malformed message: %w", err)
		}
		// Here we don't append to messages as we don't want to use
		// loads of memory. `live-update="append"` handles the appending
		// of messages in the DOM.
		m.Messages = []Message{msg}
		return m, nil
	})

	return h
}
