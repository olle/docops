package main

import "os"

func commandFor(cfg config) command {
	switch {
	case cfg.server:
		return serveCommand{
			addr:      cfg.addr,
			info:      documentSource(cfg.name),
			presenter: htmlWriter{},
			server:    stdHTTPServer{},
		}
	case cfg.tui:
		return tuiCommand{
			info:      documentSource(cfg.name),
			presenter: tuiScreen{out: os.Stdout, in: os.Stdin, term: stdioTerminal{}},
		}
	default:
		return greetCommand{
			info:      fixedInformation{msg: newMessage(cfg.name, nil)},
			presenter: plainTextWriter{},
			out:       os.Stdout,
		}
	}
}

func documentSource(name string) informationSource {
	return documentInformation{
		name:   name,
		root:   ".",
		finder: globDocumentFinder{},
		parser: fileDocumentParser{},
	}
}
