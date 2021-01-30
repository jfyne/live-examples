package components

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jfyne/live"
	"github.com/jfyne/live/component"
)

const (
	tick = "tick"
)

type ClockState struct {
	TZ   string
	Time time.Time
	loc  *time.Location
}

func (c ClockState) FormattedTime() string {
	return c.Time.Format("15:04:05")
}

func (c *ClockState) Update() {
	c.Time = time.Now().In(c.loc)
}

func NewClockState(timezone string) (*ClockState, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	now := time.Now().In(location)
	c := &ClockState{
		Time: now,
		loc:  location,
		TZ:   timezone,
	}
	return c, nil
}

func NewClock(ID string, h *live.Handler, s *live.Socket, timezone string) (component.Component, error) {
	return component.New(
		ID,
		h,
		s,
		component.WithRegister(func(c *component.Component) error {
			c.HandleSelf(tick, func(_ map[string]interface{}) (interface{}, error) {
				clock, ok := c.State.(*ClockState)
				if !ok {
					return nil, fmt.Errorf("no clock data")
				}
				clock.Update()

				go func(sock *live.Socket) {
					time.Sleep(1 * time.Second)
					c.Self(sock, live.Event{T: tick})
				}(s)

				return clock, nil
			})
			return nil
		}),
		component.WithMount(func(ctx context.Context, c *component.Component, r *http.Request, connected bool) error {
			if connected {
				go func() {
					time.Sleep(1 * time.Second)
					c.Self(s, live.Event{T: tick})
				}()
			}
			state, err := NewClockState(timezone)
			if err != nil {
				return err
			}
			c.State = state
			return nil
		}),
		component.WithRender(func(w io.Writer, c *component.Component) error {
			return c.HTML(`
                <div>
                    <p>{{.TZ}}</p>
                    <time>{{.FormattedTime}}</time>
                </div>
            `).Render(w)
		}),
	)
}
