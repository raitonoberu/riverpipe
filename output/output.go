package output

import (
	"encoding/json"
	"io"
	"github.com/raitonoberu/riverpipe/client/event"
)

type Output interface {
	Run(<-chan event.Event) error
}

type eventMessage struct {
	Event string      `json:"event"`
	Body  event.Event `json:"args"`
}

func Write(event event.Event, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(
		eventMessage{event.Event(), event},
	)
}
