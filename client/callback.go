package client

import (
	"github.com/raitonoberu/riverpipe/client/event"
	"github.com/raitonoberu/riverpipe/client/river"
)

func (c *Client) registerCallbacks() {
	c.outputStatus.SetFocusedTagsHandler(c.focusedTagsHandler)
	c.outputStatus.SetViewTagsHandler(c.viewTagsHandler)
	c.outputStatus.SetUrgentTagsHandler(c.urgentTagsHandler)
	c.outputStatus.SetLayoutNameHandler(c.layoutNameHandler)
	c.outputStatus.SetLayoutNameClearHandler(c.layoutNameClearHandler)

	c.seatStatus.SetFocusedOutputHandler(c.focusedOutputHandler)
	c.seatStatus.SetUnfocusedOutputHandler(c.unfocusedOutputHandler)
	c.seatStatus.SetFocusedViewHandler(c.focusedViewHandler)
	c.seatStatus.SetModeHandler(c.modeHandler)
}

func (c *Client) focusedTagsHandler(e river.OutputStatusFocusedTagsEvent) {
	c.putEvent(event.FocusedTags{
		Tags: e.Tags,
	})
}

func (c *Client) viewTagsHandler(e river.OutputStatusViewTagsEvent) {
	tags := make([]uint32, len(e.Tags)/4)
	for i, t := range e.Tags {
		tags[i/4] += uint32(t) << (8 * (i % 4))
	}

	c.putEvent(event.ViewTags{
		Tags: tags,
	})
}

func (c *Client) urgentTagsHandler(e river.OutputStatusUrgentTagsEvent) {
	c.putEvent(event.UrgentTags{
		Tags: e.Tags,
	})
}

func (c *Client) layoutNameHandler(e river.OutputStatusLayoutNameEvent) {
	c.putEvent(event.LayoutName{
		Name: e.Name,
	})
}

func (c *Client) layoutNameClearHandler(river.OutputStatusLayoutNameClearEvent) {
	c.putEvent(event.LayoutNameClear{})
}

func (c *Client) focusedOutputHandler(e river.SeatStatusFocusedOutputEvent) {
	c.putEvent(event.FocusedOutput{
		Output: e.Output.ID(),
	})
}

func (c *Client) unfocusedOutputHandler(e river.SeatStatusUnfocusedOutputEvent) {
	c.putEvent(event.UnfocusedOutput{
		Output: e.Output.ID(),
	})
}

func (c *Client) focusedViewHandler(e river.SeatStatusFocusedViewEvent) {
	c.putEvent(event.FocusedView{
		Title: e.Title,
	})
}

func (c *Client) modeHandler(e river.SeatStatusModeEvent) {
	c.putEvent(event.Mode{
		Name: e.Name,
	})
}
