package main

import (
	"os"

	"github.com/olle/docops/internal/app"
	"github.com/olle/docops/internal/ui/text"
	"github.com/olle/docops/internal/ui/tui"
)

func commandFor(cfg config) (command, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	switch {
	case cfg.debug:
		return debugCommand{
			source: app.NewSource(root),
			out:    os.Stdout,
		}, nil
	case cfg.server:
		return serveCommand{
			addr:      cfg.addr,
			info:      app.NewSource(root),
			presenter: text.HTML{},
			server:    stdHTTPServer{},
		}, nil
	case cfg.tui:
		return tuiCommand{
			screen: tui.New(app.NewSource(root), tui.Options{In: os.Stdin, Out: os.Stdout}),
		}, nil
	case cfg.html || cfg.markdown:
		return staticCommand{
			info:      app.NewSource(root),
			presenter: staticPresenter(cfg.markdown),
			out:       os.Stdout,
		}, nil
	default:
		return greetCommand{
			info:      app.NewFixedSource(defaultName, nil),
			presenter: text.PlainText{},
			out:       os.Stdout,
		}, nil
	}
}

func staticPresenter(markdown bool) app.Presenter {
	if markdown {
		return text.Markdown{}
	}
	return text.HTML{}
}
