package components

import (
	"context"
	"fmt"

	"github.com/jfyne/live"
	"github.com/jfyne/live/page"
	g "github.com/maragudk/gomponents"
	co "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
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
	return co.HTML5(co.HTML5Props{
		Title:    c.Title,
		Language: "en",
		Head: []g.Node{
			h.StyleEl(h.Type("text/css"),
				g.Raw(`body {font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol"; }`),
			),
		},
		Body: []g.Node{
			h.H1(g.Text("World Clocks")),
			h.FormEl(
				h.ID("tz-form"),
				g.Attr("live-change", c.Event("validate-tz")), // c.Event scopes the events to this component.
				g.Attr("live-submit", c.Event("add-time")),
				h.Div(
					h.P(g.Text("Try Europe/London or America/New_York")),
					h.Input(h.Name("tz")),
					g.If(c.ValidationError != "", h.Span(g.Text(c.ValidationError))),
				),
				h.Input(h.Type("submit"), g.If(c.ValidationError != "", h.Disabled())),
			),
			h.Div(
				g.Group(g.Map(len(c.Clocks), func(idx int) g.Node {
					return c.Clocks[idx].Render()
				})),
			),
			h.Script(h.Src("/live.js")),
		},
	}).Render
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
	//
	if err := page.Start(ctx, fmt.Sprintf("clock-%d", len(c.Clocks)+1), c.Handler, c.Socket, clock); err != nil {
		return err
	}

	// Update the page state with the new clock.
	c.Clocks = append(c.Clocks, clock)

	return nil
}