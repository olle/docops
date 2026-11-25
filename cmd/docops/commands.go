package main

import (
	"io"
	"log"
)

type command interface {
	Run() error
}

func messageFor(info informationSource) message {
	msg, err := info.Message()
	if err != nil {
		log.Printf("reading information: %v", err)
	}
	return msg
}

type greetCommand struct {
	info      informationSource
	presenter presenter
	out       io.Writer
}

func (c greetCommand) Run() error {
	return c.presenter.Render(c.out, messageFor(c.info))
}

type serveCommand struct {
	addr      string
	info      informationSource
	presenter webPresenter
	server    httpServer
}

func (c serveCommand) Run() error {
	return c.server.ListenAndServe(c.addr, handlerFor(c.presenter, messageFor(c.info)))
}

type tuiCommand struct {
	info      informationSource
	presenter screenPresenter
}

func (c tuiCommand) Run() error {
	return c.presenter.Present(messageFor(c.info))
}
