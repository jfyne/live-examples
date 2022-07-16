package components

import (
	"context"
	"time"

	"github.com/jfyne/live/page"
)

// Clock is a timezone aware clock.
type Clock struct {
	TZ   string
	Time time.Time
	loc  *time.Location

	page.Component
}

// NewClock create a new clock component.
func NewClock(timezone string) (*Clock, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	now := time.Now().In(location)
	c := &Clock{
		Time: now,
		loc:  location,
		TZ:   timezone,
	}
	return c, nil
}

func (c Clock) FormattedTime() string {
	return c.Time.Format(time.RubyDate)
}

func (c *Clock) Mount(ctx context.Context) error {
	// If we are mounting on connection send the first tick event.
	if c.Socket.Connected() {
		go func() {
			time.Sleep(1 * time.Second)
			c.Self(ctx, "tick", time.Now())
		}()
	}
	return nil
}

func (c Clock) Render() page.RenderFunc {
	return page.HTML(`
        <div>
            <p>{{.TZ}}</p>
            <time datetime="{{.Time}}">{{.FormattedTime}}</time>
        </div>
    `, c)
}

func (c *Clock) OnTick(ctx context.Context, t time.Time) {
	c.Time = t.In(c.loc)
	go func() {
		time.Sleep(1 * time.Second)
		c.Self(ctx, "tick", time.Now())
	}()
}
