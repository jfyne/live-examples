package components

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live/page"
	"github.com/maragudk/gomponents"
	comp "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/html"
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
	c.HandleEvent(validateTZ, func(ctx context.Context, p map[string]interface{}) (interface{}, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the tz coming from the form.
		tz := live.ParamString(p, "tz")

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
	c.HandleEvent(addTime, func(ctx context.Context, p map[string]interface{}) (interface{}, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the timezone sent from the form input.
		tz := live.ParamString(p, "tz")
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
func pageRender(w io.Writer, c *page.Component) error {
	state, ok := c.State.(*PageState)
	if !ok {
		return fmt.Errorf("could not get state")
	}

	// Here we use the gomponents library to do typed rendering.
	// https://github.com/maragudk/gomponents
	return comp.HTML5(comp.HTML5Props{
		Title:    state.Title,
		Language: "en",
		Head: []gomponents.Node{
			html.StyleEl(html.Type("text/css"),
				gomponents.Raw(`body {font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol"; }`),
			),
		},
		Body: []gomponents.Node{
			html.H1("World Clocks"),
			html.FormEl(
				html.ID("tz-form"),
				gomponents.Attr("live-change", c.Event(validateTZ)), // c.Event scopes the events to this component.
				gomponents.Attr("live-submit", c.Event(addTime)),
				html.Div(
					html.P(gomponents.Text("Try Europe/London or America/New_York")),
					html.Input(html.Name("tz")),
					gomponents.If(state.ValidationError != "", html.Span(gomponents.Text(state.ValidationError))),
				),
				html.Input(html.Type("submit"), gomponents.If(state.ValidationError != "", html.Disabled())),
			),
			html.Div(
				gomponents.Group(gomponents.Map(len(state.Timezones), func(idx int) gomponents.Node {
					return page.Render(state.Timezones[idx])
				})),
			),
			html.Script(html.Src("/live.js")),
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
