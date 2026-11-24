package main

import (
	"log"
	"os"
)

type command interface {
	Run() error
}

type greetCommand struct {
	name   string
	writer outputWriter
}

func (c greetCommand) Run() error {
	return c.writer.Render(os.Stdout, newMessage(c.name, nil))
}

type serveCommand struct {
	addr   string
	root   string
	name   string
	finder documentFinder
	parser documentParser
	writer outputWriter
	server httpServer
}

func (c serveCommand) Run() error {
	titles, err := collectTitles(c.root, c.finder, c.parser)
	if err != nil {
		log.Printf("reading AsyncAPI documents: %v", err)
	}
	return c.server.ListenAndServe(c.addr, handlerFor(c.writer, newMessage(c.name, titles)))
}
