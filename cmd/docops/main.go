package main

import (
	"flag"
	"log"
)

func main() {
	name := flag.String("name", "World", "name to greet")
	serverMode := flag.Bool("server", false, "run an HTTP server that greets visitors")
	addr := flag.String("addr", defaultAddr, "address to listen on in server mode")
	flag.Parse()

	textWriter := plainTextWriter{}
	htmlOutput := htmlWriter{}
	var cmd command
	if *serverMode {
		cmd = serveCommand{
			addr:   *addr,
			root:   ".",
			name:   *name,
			finder: globDocumentFinder{},
			parser: fileDocumentParser{},
			writer: htmlOutput,
			server: stdHTTPServer{},
		}
	} else {
		cmd = greetCommand{name: *name, writer: textWriter}
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

const defaultAddr = ":8080"
