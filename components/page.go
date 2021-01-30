package components

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jfyne/live"
	"github.com/jfyne/live/component"
	"github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/html"
)

const (
	validateTZ = "validate-tz"
	addTime    = "add-time"
)

type PageState struct {
	Title           string
	ValidationError string
	Timezones       []component.Component
}

func pageState(title string) *PageState {
	return &PageState{
		Title:     title,
		Timezones: []component.Component{},
	}
}

func NewPage(ID string, h *live.Handler, s *live.Socket, title string) (component.Component, error) {
	return component.New(
		ID,
		h,
		s,
		component.WithRegister(func(c *component.Component) error {
			// Handler for the timezone entry validation.
			c.HandleEvent(validateTZ, func(p map[string]interface{}) (interface{}, error) {
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
			c.HandleEvent(addTime, func(p map[string]interface{}) (interface{}, error) {
				// Get the current page component state.
				state, _ := c.State.(*PageState)

				// Get the timezone sent from the form input.
				tz := live.ParamString(p, "tz")
				if tz == "" {
					return state, nil
				}

				clock, err := component.Init(context.Background(), func() (component.Component, error) {
					return NewClock(fmt.Sprintf("clock-%d", len(state.Timezones)+1), c.Handler, c.Socket, tz)
				})
				if err != nil {
					return state, err
				}

				state.Timezones = append(state.Timezones, clock)
				return state, nil
			})

			return nil
		}),
		component.WithMount(func(ctx context.Context, c *component.Component, r *http.Request, connected bool) error {
			c.State = pageState("Clocks")
			return nil
		}),
		component.WithRender(func(w io.Writer, c *component.Component) error {
			state, ok := c.State.(*PageState)
			if !ok {
				return fmt.Errorf("could not get state")
			}

			return components.HTML5(components.HTML5Props{
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
						gomponents.Attr("live-change", c.EventPrefix(validateTZ)),
						gomponents.Attr("live-submit", c.EventPrefix(addTime)),
						html.Div(
							html.P(gomponents.Text("Try Europe/London or America/New_York")),
							html.Input(html.Name("tz")),
							gomponents.If(state.ValidationError != "", html.Span(gomponents.Text(state.ValidationError))),
						),
						html.Input(html.Type("submit"), gomponents.If(state.ValidationError != "", html.Disabled())),
					),
					html.Div(
						gomponents.Group(gomponents.Map(len(state.Timezones), func(idx int) gomponents.Node {
							return component.RenderComponent(state.Timezones[idx])
						})),
					),
					html.Script(html.Src("/live.js")),
				},
			}).Render(w)
		}),
	)
}
