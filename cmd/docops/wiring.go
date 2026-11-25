package main

import "os"

func commandFor(cfg config) (command, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	switch {
	case cfg.server:
		return serveCommand{
			addr:      cfg.addr,
			info:      documentSource(root),
			presenter: htmlWriter{},
			server:    stdHTTPServer{},
		}, nil
	case cfg.tui:
		return tuiCommand{
			info:      documentSource(root),
			presenter: tuiScreen{out: os.Stdout, in: os.Stdin, term: stdioTerminal{}},
		}, nil
	case cfg.html || cfg.markdown:
		return staticCommand{
			info:      documentSource(root),
			presenter: staticPresenter(cfg.markdown),
			out:       os.Stdout,
		}, nil
	default:
		return greetCommand{
			info:      fixedInformation{msg: newMessage(defaultName, nil)},
			presenter: plainTextWriter{},
			out:       os.Stdout,
		}, nil
	}
}

func documentSource(root string) informationSource {
	return documentInformation{
		name:   defaultName,
		root:   root,
		finder: documentNameFinder{},
		parser: fileDocumentParser{},
	}
}

func staticPresenter(markdown bool) presenter {
	if markdown {
		return markdownWriter{}
	}
	return htmlWriter{}
}
