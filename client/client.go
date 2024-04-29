package client

import (
	"errors"
	"fmt"
	"github.com/raitonoberu/riverpipe/client/event"
	"github.com/raitonoberu/riverpipe/client/river"

	"github.com/rajveermalviya/go-wayland/wayland/client"
)

func New() (*Client, error) {
	client := &Client{}
	if err := client.initInterfaces(); err != nil {
		return nil, err
	}
	if err := client.initRiverClients(); err != nil {
		return nil, err
	}
	client.registerCallbacks()
	return client, nil
}

type Client struct {
	display  *client.Display
	registry *client.Registry
	output   *client.Output
	seat     *client.Seat

	statusManager *river.StatusManager
	outputStatus  *river.OutputStatus
	seatStatus    *river.SeatStatus

	eventCh chan<- event.Event
	err     error
}

func (c *Client) Run(ch chan<- event.Event) error {
	c.eventCh = ch

	for {
		if c.err != nil {
			return c.err
		}

		if err := c.dispatch(); err != nil {
			return err
		}
	}
}

func (c *Client) Release() error {
	if c.seatStatus != nil {
		if err := c.seatStatus.Destroy(); err != nil {
			return err
		}
		c.seatStatus = nil
	}
	if c.outputStatus != nil {
		if err := c.outputStatus.Destroy(); err != nil {
			return err
		}
		c.outputStatus = nil
	}
	if c.statusManager != nil {
		if err := c.statusManager.Destroy(); err != nil {
			return err
		}
		c.statusManager = nil
	}

	if c.seat != nil {
		if err := c.seat.Release(); err != nil {
			return err
		}
		c.seat = nil
	}
	if c.output != nil {
		if err := c.output.Release(); err != nil {
			return err
		}
		c.output = nil
	}
	if c.registry != nil {
		if err := c.registry.Destroy(); err != nil {
			return err
		}
		c.registry = nil
	}
	if c.display != nil {
		if err := c.display.Destroy(); err != nil {
			return err
		}
		c.display = nil
	}

	return nil
}

func (c *Client) initInterfaces() error {
	display, err := client.Connect("")
	if err != nil {
		return fmt.Errorf("failed to connect to display: %w", err)
	}

	display.SetErrorHandler(c.handleDisplayError)
	c.display = display

	registry, err := c.display.GetRegistry()
	if err != nil {
		return fmt.Errorf("failed to get global registry object: %w", err)
	}
	registry.SetGlobalHandler(c.handleRegistryGlobal)
	c.registry = registry

	if err := c.syncDisplay(); err != nil {
		return err
	}
	if c.statusManager == nil {
		return errors.New("couldn't connect to zriver_status_manager_v1. Is River running?")
	}
	return nil
}

func (c *Client) initRiverClients() error {
	outputStatus, err := c.statusManager.GetRiverOutputStatus(c.output)
	if err != nil {
		return fmt.Errorf("failed to get river output status: %w", err)
	}
	c.outputStatus = outputStatus

	seatStatus, err := c.statusManager.GetRiverSeatStatus(c.seat)
	if err != nil {
		return fmt.Errorf("failed to get river sear status: %w", err)
	}
	c.seatStatus = seatStatus

	return nil
}

func (c *Client) syncDisplay() error {
	callback, err := c.display.Sync()
	if err != nil {
		return fmt.Errorf("failed to get sync callback: %w", err)
	}
	defer callback.Destroy()

	done := false
	callback.SetDoneHandler(func(_ client.CallbackDoneEvent) {
		done = true
	})

	for !done {
		if err := c.dispatch(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) handleRegistryGlobal(e client.RegistryGlobalEvent) {
	switch e.Interface {
	case "wl_output":
		output := client.NewOutput(c.context())
		err := c.registry.Bind(e.Name, e.Interface, e.Version, output)
		if err != nil {
			c.err = fmt.Errorf("failed to bind wl_output interface: %w", err)
		}
		c.output = output
	case "wl_seat":
		if c.seat != nil {
			return
		}
		seat := client.NewSeat(c.context())
		err := c.registry.Bind(e.Name, e.Interface, e.Version, seat)
		if err != nil {
			c.err = fmt.Errorf("failed to bind wl_seat interface: %w", err)
		}
		c.seat = seat
	case "zriver_status_manager_v1":
		sm := river.NewStatusManager(c.context())
		err := c.registry.Bind(e.Name, e.Interface, e.Version, sm)
		if err != nil {
			c.err = fmt.Errorf("failed to bind zriver_status_manager_v1 interface: %w", err)
		}
		c.statusManager = sm
	}
}

func (c *Client) putEvent(event event.Event) {
	if c.eventCh == nil {
		return
	}
	c.eventCh <- event
}

func (c *Client) handleDisplayError(event client.DisplayErrorEvent) {
	c.err = fmt.Errorf("display error: %s (code %d)", event.Message, event.Code)
}
