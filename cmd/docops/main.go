package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const defaultAddr = ":8080"

func main() {
	name := flag.String("name", "World", "name to greet")
	serverMode := flag.Bool("server", false, "run an HTTP server that greets visitors")
	addr := flag.String("addr", defaultAddr, "address to listen on in server mode")
	flag.Parse()

	if *serverMode {
		if err := newHelloServer(*addr).ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return
	}

	fmt.Printf("Hello %s!\n", *name)
}
