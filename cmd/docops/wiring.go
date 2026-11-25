package main

import "os"

func commandFor(cfg config) command {
	switch {
	case cfg.server:
		return serveCommand{
			addr:      cfg.addr,
			info:      documentSource(),
			presenter: htmlWriter{},
			server:    stdHTTPServer{},
		}
	case cfg.tui:
		return tuiCommand{
			info:      documentSource(),
			presenter: tuiScreen{out: os.Stdout, in: os.Stdin, term: stdioTerminal{}},
		}
	case cfg.html || cfg.markdown:
		return staticCommand{
			info:      documentSource(),
			presenter: staticPresenter(cfg.markdown),
			out:       os.Stdout,
		}
	default:
		return greetCommand{
			info:      fixedInformation{msg: newMessage(defaultName, nil)},
			presenter: plainTextWriter{},
			out:       os.Stdout,
		}
	}
}

func documentSource() informationSource {
	return documentInformation{
		name:   defaultName,
		root:   ".",
		finder: globDocumentFinder{},
		parser: fileDocumentParser{},
	}
}

func staticPresenter(markdown bool) presenter {
	if markdown {
		return markdownWriter{}
	}
	return htmlWriter{}
}
