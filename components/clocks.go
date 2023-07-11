package components

import (
	"context"
	"fmt"

	"github.com/jfyne/live"
	"github.com/jfyne/live/page"
)

// Clocks component handles creating new timezone clocks.
type Clocks struct {
	Title           string
	ValidationError string
	Clocks          []*Clock

	page.Component
}

// NewClocks create a new clocks component.
func NewClocks(title string) (*Clocks, error) {
	return &Clocks{
		Title:  title,
		Clocks: []*Clock{},
	}, nil
}

func (c Clocks) Render() page.RenderFunc {
	return page.HTML(`
    <!doctype html>
    <html lang="en">
      <head>
        <meta charset="utf-8" />
        <meta http-equiv="x-ua-compatible" content="ie=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>{{.Title}}</title>
        <style>
            body {font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol"; }
        </style>
      </head>
      <body>
        <h1>World Clocks</h1>
        <form id="tz-form" live-change="{{ Event "validate-tz" }}" live-submit="{{ Event "add-time" }}">
            <div>
                <p>Try Europe/London or America/New_York</p>
                <input name="tz">
                {{ if ne .ValidationError "" }}
                    <span>{{ .ValidationError }}</span>
                {{ end }}
            </div>
            <input type="submit" {{ if ne .ValidationError "" }} diabaled="true" {{ end }}>
        </form>
        <div>
            {{ range .Clocks }}
                {{ Component . }}
            {{ end }}
        </div>
        <script src="/live.js"></script>
      </body>
    </html>
    `, c)
}

func (c *Clocks) OnValidateTz(ctx context.Context, p live.Params) error {
	// Get the tz coming from the form.
	tz := p.String("tz")

	// Try to make a new Clock, this will return an error if the
	// timezone is not real.
	if _, err := NewClock(tz); err != nil {
		c.ValidationError = fmt.Sprintf("Timezone %s does not exist", tz)
		return nil
	}

	// If there was no error loading the clock state reset the
	// validation error.
	c.ValidationError = ""

	return nil
}

func (c *Clocks) OnAddTime(ctx context.Context, p live.Params) error {
	// Get the timezone sent from the form input.
	tz := p.String("tz")
	if tz == "" {
		return nil
	}

	// Create the new clock with the timezone.
	clock, err := NewClock(tz)
	if err != nil {
		return err
	}
	if err := page.Start(ctx, fmt.Sprintf("clock-%d", len(c.Clocks)+1), c.Handler, c.Socket, clock); err != nil {
		return err
	}

	// Update the page state with the new clock.
	c.Clocks = append(c.Clocks, clock)

	return nil
}
