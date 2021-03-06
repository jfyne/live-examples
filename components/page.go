package components

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live/page"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

const (
	validateTZ = "validate-tz"
	addTime    = "add-time"
)

// PageState the state we are tracking for our page.
type PageState struct {
	Title           string
	ValidationError string
	Timezones       []*page.Component
}

// newPageState create a new page state.
func newPageState(title string) *PageState {
	return &PageState{
		Title:     title,
		Timezones: []*page.Component{},
	}
}

// pageRegister register the pages events.
func pageRegister(c *page.Component) error {
	// Handler for the timezone entry validation.
	c.HandleEvent(validateTZ, func(ctx context.Context, p live.Params) (interface{}, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the tz coming from the form.
		tz := p.String("tz")

		// Try to make a new ClockState, this will return an error if the
		// timezone is not real.
		if _, err := NewClockState(tz); err != nil {
			state.ValidationError = fmt.Sprintf("Timezone %s does not exist", tz)
			return state, nil
		}

		// If there was no error loading the clock state reset the
		// validation error.
		state.ValidationError = ""

		return state, nil
	})

	// Handler for adding a timezone.
	c.HandleEvent(addTime, func(ctx context.Context, p live.Params) (interface{}, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the timezone sent from the form input.
		tz := p.String("tz")
		if tz == "" {
			return state, nil
		}

		// Use the page.Init function to create a new clock, register it and mount it.
		clock, err := page.Init(context.Background(), func() (*page.Component, error) {
			// Each clock requires its own unique stable ID. Events for each clock can then find
			// their own component.
			return NewClock(fmt.Sprintf("clock-%d", len(state.Timezones)+1), c.Handler, c.Socket, tz)
		})
		if err != nil {
			return state, err
		}

		// Update the page state with the new clock.
		state.Timezones = append(state.Timezones, clock)

		// Return the state to have it persisted.
		return state, nil
	})

	return nil
}

// pageMount initialise the page component.
func pageMount(title string) page.MountHandler {
	return func(ctx context.Context, c *page.Component, r *http.Request) error {
		// Create a new page state.
		c.State = newPageState(title)
		return nil
	}
}

// pageRender render the page component.
func pageRender(w io.Writer, cmp *page.Component) error {
	state, ok := cmp.State.(*PageState)
	if !ok {
		return fmt.Errorf("could not get state")
	}

	// Here we use the gomponents library to do typed rendering.
	// https://github.com/maragudk/gomponents
	return c.HTML5(c.HTML5Props{
		Title:    state.Title,
		Language: "en",
		Head: []g.Node{
			StyleEl(Type("text/css"),
				g.Raw(`body {font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol"; }`),
			),
		},
		Body: []g.Node{
			H1(g.Text("World Clocks")),
			FormEl(
				ID("tz-form"),
				g.Attr("live-change", cmp.Event(validateTZ)), // c.Event scopes the events to this component.
				g.Attr("live-submit", cmp.Event(addTime)),
				Div(
					P(g.Text("Try Europe/London or America/New_York")),
					Input(Name("tz")),
					g.If(state.ValidationError != "", Span(g.Text(state.ValidationError))),
				),
				Input(Type("submit"), g.If(state.ValidationError != "", Disabled())),
			),
			Div(
				g.Group(g.Map(len(state.Timezones), func(idx int) g.Node {
					return page.Render(state.Timezones[idx])
				})),
			),
			Script(Src("/live.js")),
		},
	}).Render(w)
}

// NewPage create a new page component.
func NewPage(ID string, h *live.Handler, s *live.Socket, title string) (*page.Component, error) {
	return page.NewComponent(ID, h, s,
		page.WithRegister(pageRegister),
		page.WithMount(pageMount(title)),
		page.WithRender(pageRender),
	)
}
