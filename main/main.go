package main

import (
	"ares"
	"flag"
	"log"
)

var tls = flag.Bool("tls", true, "whether TLS is used")
var port = flag.Int("port", 443, "The TCP port that the server listen on")

func main() {
	flag.Parse()

	log.Printf("Starting ares: tls = %t, port = %d", *tls, *port)

	ares := AresController.NewAresController()
	ares.Run(*port, *tls)
}
