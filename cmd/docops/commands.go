package main

import (
	"context"
	"io"
	"log"

	"github.com/olle/docops/internal/app"
	"github.com/olle/docops/internal/ui/text"
	"github.com/olle/docops/internal/view"
)

type command interface {
	Run() error
}

func snapshotFor(source app.SnapshotSource) view.Snapshot {
	snapshot, err := source.Snapshot(context.Background())
	if err != nil {
		log.Printf("reading information: %v", err)
	}
	return snapshot
}

type debugCommand struct {
	source app.SnapshotSource
	out    io.Writer
}

func (c debugCommand) Run() error {
	snapshot, err := c.source.Snapshot(context.Background())
	if err != nil {
		return err
	}
	return (text.Debug{}).Render(c.out, snapshot)
}

type greetCommand struct {
	info      app.SnapshotSource
	presenter app.Presenter
	out       io.Writer
}

func (c greetCommand) Run() error {
	return c.presenter.Render(c.out, snapshotFor(c.info))
}

type staticCommand struct {
	info      app.SnapshotSource
	presenter app.Presenter
	out       io.Writer
}

func (c staticCommand) Run() error {
	return c.presenter.Render(c.out, snapshotFor(c.info))
}

type serveCommand struct {
	addr      string
	info      app.SnapshotSource
	presenter app.WebPresenter
	server    httpServer
}

func (c serveCommand) Run() error {
	return c.server.ListenAndServe(c.addr, handlerFor(c.presenter, snapshotFor(c.info)))
}

type tuiCommand struct {
	screen app.Screen
}

func (c tuiCommand) Run() error {
	return c.screen.Run(context.Background())
}
