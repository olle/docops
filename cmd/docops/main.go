package main

import (
	"errors"
	"flag"
	"log"
	"os"

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
	d := fs.Bool("d", false, "overview debug output")
	dd := fs.Bool("dd", false, "medium debug output")
	ddd := fs.Bool("ddd", false, "full debug output")
	debug := fs.Bool("debug", false, "alias for -d")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if fs.NArg() > 0 {
		return config{}, errors.New("unexpected positional arguments")
	}

	debugDepth := 0
	for _, flag := range []struct {
		enabled *bool
		depth   int
	}{
		{enabled: d, depth: 1},
		{enabled: dd, depth: 2},
		{enabled: ddd, depth: 3},
		{enabled: debug, depth: 1},
	} {
		if !*flag.enabled {
			continue
		}
		if debugDepth > 0 {
			return config{}, errors.New("only one debug depth flag can be used at a time")
		}
		debugDepth = flag.depth
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
