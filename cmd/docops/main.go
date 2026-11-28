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
	web      bool
	tui      bool
	html     bool
	markdown bool
	txt      bool
	pdf      bool
	debug    int
	addr     string
}

type webFlag struct {
	enabled bool
	addr    string
}

func (f *webFlag) String() string {
	if f.addr == "" {
		return defaultAddr
	}
	return f.addr
}

func (f *webFlag) Set(value string) error {
	switch value {
	case "false":
		f.enabled = false
		return nil
	case "", "true":
		value = defaultAddr
	}
	f.enabled = true
	f.addr = value
	return nil
}

func (*webFlag) IsBoolFlag() bool {
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
	fs := flag.NewFlagSet("docops", flag.ContinueOnError)

	web := webFlag{addr: defaultAddr}
	fs.Var(&web, "web", "run a web server, optionally at -web=ADDR")
	tui := fs.Bool("tui", false, "run the interactive terminal interface")
	html := fs.Bool("html", false, "render static HTML")
	md := fs.Bool("md", false, "render static Markdown")
	txt := fs.Bool("txt", false, "render static plain text")
	pdf := fs.Bool("pdf", false, "render a published PDF catalog")

	visibleArgs, debugDepth, err := parseDebugArgs(args)
	if err != nil {
		return config{}, err
	}
	if err := fs.Parse(visibleArgs); err != nil {
		return config{}, err
	}
	if fs.NArg() > 0 {
		return config{}, errors.New("unexpected positional arguments")
	}

	modes := 0
	for _, enabled := range []bool{web.enabled, *tui, *html, *md, *txt, *pdf, debugDepth > 0} {
		if enabled {
			modes++
		}
	}
	if modes > 1 {
		return config{}, errors.New("only one of -html, -md, -txt, -pdf, -web, -tui, -d, -dd, -ddd, or -debug can be used at a time")
	}

	return config{
		web:      web.enabled,
		tui:      *tui,
		html:     *html,
		markdown: *md,
		txt:      *txt,
		pdf:      *pdf,
		debug:    debugDepth,
		addr:     web.addr,
	}, nil
}

func parseDebugArgs(args []string) ([]string, int, error) {
	depths := map[string]int{"d": 1, "dd": 2, "ddd": 3, "debug": 1}
	values := map[string]bool{}
	remaining := make([]string, 0, len(args))

	for _, arg := range args {
		name, value, hasValue := strings.Cut(arg, "=")
		name = strings.TrimPrefix(strings.TrimPrefix(name, "-"), "-")
		if _, ok := depths[name]; !ok {
			remaining = append(remaining, arg)
			continue
		}

		enabled := true
		if hasValue {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid boolean value %q for -%s: %v", value, name, err)
			}
			enabled = parsed
		}
		values[name] = enabled
	}

	depth := 0
	for _, name := range []string{"d", "dd", "ddd", "debug"} {
		if !values[name] {
			continue
		}
		if depth > 0 {
			return nil, 0, errors.New("only one debug depth flag can be used at a time")
		}
		depth = depths[name]
	}
	return remaining, depth, nil
}
