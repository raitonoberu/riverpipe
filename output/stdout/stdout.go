package stdout

import (
	"os"
	"github.com/raitonoberu/riverpipe/client/event"
	"github.com/raitonoberu/riverpipe/output"
)

func New() *Stdout {
	return &Stdout{}
}

type Stdout struct{}

func (s *Stdout) Run(ch <-chan event.Event) error {
	for e := range ch {
		err := output.Write(e, os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}
