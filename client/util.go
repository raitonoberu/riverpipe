package client

import "github.com/rajveermalviya/go-wayland/wayland/client"

func (c *Client) context() *client.Context {
	return c.display.Context()
}

func (c *Client) dispatch() error {
	return c.context().Dispatch()
}
