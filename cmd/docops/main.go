package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olle/docops/internal/app"
)

const (
	defaultAddr = "127.0.0.1:8080"
	defaultName = app.DefaultName
)

type config struct {
	server   bool
	tui      bool
	html     bool
	markdown bool
	debug    bool
	addr     string
}

type serverFlag struct {
	enabled bool
	addr    string
}

func (f *serverFlag) String() string {
	if f.addr == "" {
		return defaultAddr
	}
	return f.addr
}

func (f *serverFlag) Set(value string) error {
	if value == "false" {
		f.enabled = false
		return nil
	}
	if value == "" || value == "true" {
		value = defaultAddr
	}
	f.enabled = true
	f.addr = value
	return nil
}

func (*serverFlag) IsBoolFlag() bool {
	return true
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	cmd, err := commandFor(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(args []string) (config, error) {
	args, debug, err := parseDebugArgs(args)
	if err != nil {
		return config{}, err
	}

	fs := flag.NewFlagSet("docops", flag.ContinueOnError)

	server := serverFlag{addr: defaultAddr}
	fs.Var(&server, "server", "run an HTTP server, optionally at -server=ADDR")
	tui := fs.Bool("tui", false, "run the interactive terminal interface")
	html := fs.Bool("html", false, "render static HTML")
	md := fs.Bool("md", false, "render static Markdown")

	if err := fs.Parse(normalizeServerArgs(args)); err != nil {
		return config{}, err
	}

	modes := 0
	for _, enabled := range []bool{server.enabled, *tui, *html, *md, debug} {
		if enabled {
			modes++
		}
	}
	if modes > 1 {
		return config{}, errors.New("only one of -server, -tui, -html, -md, or -debug can be used at a time")
	}

	return config{
		server:   server.enabled,
		tui:      *tui,
		html:     *html,
		markdown: *md,
		debug:    debug,
		addr:     server.addr,
	}, nil
}

func normalizeServerArgs(args []string) []string {
	normalized := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "-server" || arg == "--server") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			normalized = append(normalized, arg+"="+args[i+1])
			i++
			continue
		}
		normalized = append(normalized, arg)
	}
	return normalized
}

func parseDebugArgs(args []string) ([]string, bool, error) {
	remaining := make([]string, 0, len(args))
	debug := false
	for _, arg := range args {
		name, value, hasValue := strings.Cut(arg, "=")
		switch name {
		case "-debug", "--debug", "-d", "--d":
			if !hasValue {
				debug = true
				continue
			}
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				return nil, false, fmt.Errorf("invalid debug value %q: %w", value, err)
			}
			debug = enabled
		default:
			remaining = append(remaining, arg)
		}
	}
	return remaining, debug, nil
}
