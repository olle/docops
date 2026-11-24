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

	writer := plainTextWriter{}
	var cmd command
	if *serverMode {
		cmd = serveCommand{
			addr:   *addr,
			root:   ".",
			finder: globDocumentFinder{},
			parser: fileDocumentParser{},
			writer: writer,
			server: stdHTTPServer{},
		}
	} else {
		cmd = greetCommand{name: *name, writer: writer}
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

const defaultAddr = ":8080"
