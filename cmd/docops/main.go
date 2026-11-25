package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

const defaultAddr = "127.0.0.1:8080"

type config struct {
	name   string
	server bool
	tui    bool
	addr   string
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := commandFor(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(args []string) (config, error) {
	fs := flag.NewFlagSet("docops", flag.ContinueOnError)

	name := fs.String("name", "World", "name to greet")
	server := fs.Bool("server", false, "run an HTTP server that greets visitors")
	tui := fs.Bool("tui", false, "run the interactive terminal interface")
	addr := fs.String("addr", defaultAddr, "address to listen on in server mode")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if *server && *tui {
		return config{}, errors.New("-server and -tui cannot be used together")
	}
	return config{name: *name, server: *server, tui: *tui, addr: *addr}, nil
}
