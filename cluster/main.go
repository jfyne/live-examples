package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live-examples/chat"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/mempubsub"
)

const app = "chat-app"

type CloudTransport struct {
	topic *pubsub.Topic
}

func NewCloudTransport(ctx context.Context) (*CloudTransport, error) {
	topic, err := pubsub.OpenTopic(ctx, "mem://broadcast")
	if err != nil {
		return nil, err
	}
	return &CloudTransport{
		topic: topic,
	}, nil
}

func (c *CloudTransport) Publish(ctx context.Context, topic string, msg live.Event) error {
	data, err := json.Marshal(live.TransportMessage{Topic: topic, Msg: msg})
	if err != nil {
		return fmt.Errorf("could not publish event: %w", err)
	}
	return c.topic.Send(ctx, &pubsub.Message{
		Body: data,
		Metadata: map[string]string{
			"topic": topic,
		},
	})
}

func (c *CloudTransport) Listen(ctx context.Context, p *live.PubSub) error {
	sub, err := pubsub.OpenSubscription(ctx, "mem://broadcast")
	if err != nil {
		return fmt.Errorf("could not open subscription: %w", err)
	}
	for {
		msg, err := sub.Receive(ctx)
		if err != nil {
			log.Println("receive message failed: %w", err)
			break
		}

		var t live.TransportMessage
		if err := json.Unmarshal(msg.Body, &t); err != nil {
			log.Println("malformed message received: %w", err)
			continue
		}
		p.Recieve(t.Topic, t.Msg)
		msg.Ack()
	}
	return fmt.Errorf("stopped receiving messages")
}

func main() {
	chat1 := chat.NewHandler()
	chat2 := chat.NewHandler()

	ctx := context.Background()

	t, err := NewCloudTransport(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pubsub := live.NewPubSub(ctx, t)
	pubsub.Subscribe(app, chat1)
	pubsub.Subscribe(app, chat2)

	// Run the server.
	http.Handle("/one", chat1)
	http.Handle("/two", chat2)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	http.ListenAndServe(":8080", nil)
}
