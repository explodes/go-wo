package main

import (
	"flag"

	"log"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/templates/multiscene/internal"
)

var (
	debug = flag.Bool("debug", false, "enable debug mode")
)

func main() {
	flag.Parse()

	w := internal.NewWorld(*debug)

	wo.Run(func() {
		if err := w.Run(); err != nil {
			log.Fatal(err)
		}
	})
}
