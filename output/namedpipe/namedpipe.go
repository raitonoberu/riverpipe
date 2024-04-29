package namedpipe

import (
	"os"
	"github.com/raitonoberu/riverpipe/client/event"
	"github.com/raitonoberu/riverpipe/output"
	"syscall"
)

func New(path string) *NamedPipe {
	return &NamedPipe{
		path: path,
	}
}

type NamedPipe struct {
	path string
}

func (n *NamedPipe) Run(ch <-chan event.Event) error {
	os.Remove(n.path)
	err := syscall.Mkfifo(n.path, 0666)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(n.path, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return err
	}

	for e := range ch {
		err := output.Write(e, file)
		if err != nil {
			return err
		}
	}
	return nil
}
