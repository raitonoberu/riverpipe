package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/raitonoberu/riverpipe/client"
	"github.com/raitonoberu/riverpipe/client/event"
	"github.com/raitonoberu/riverpipe/output"
	"github.com/raitonoberu/riverpipe/output/namedpipe"
	"github.com/raitonoberu/riverpipe/output/stdout"

	"github.com/urfave/cli/v2"
)

func main() {
	var file string
	var bufsize int

	app := &cli.App{
		Name:  "riverpipe",
		Usage: "print River events in JSON format",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "where to create named pipe (will be removed if exists!)",
				Destination: &file,
			},
			&cli.IntFlag{
				Name:        "bufsize",
				Aliases:     []string{"b"},
				Usage:       "events buffer size",
				Value:       16,
				Destination: &bufsize,
			},
		},
		Action: func(ctx *cli.Context) error {
			var out output.Output
			if file == "" {
				out = stdout.New()
			} else {
				out = namedpipe.New(file)
			}

			client, err := client.New()
			if err != nil {
				return fmt.Errorf("couldn't create client: %w", err)
			}
			defer client.Release()

			cleanup := make(chan os.Signal, 1)
			signal.Notify(cleanup, os.Interrupt, os.Kill)
			go func() {
				<-cleanup
				client.Release()
				os.Exit(1)
			}()

			ch := make(chan event.Event, bufsize)

			go func() {
				err := client.Run(ch)
				if err != nil {
					fmt.Fprintln(os.Stderr, "error:", err.Error())
					client.Release()
					os.Exit(1)
				}
			}()
			return out.Run(ch)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err.Error())
		os.Exit(1)
	}
}
