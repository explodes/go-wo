package main

import (
	"flag"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/behavioral/internal"
)

var (
	debug = flag.Bool("debug", false, "enable debug drawing")
)

func main() {
	flag.Parse()
	g := internal.NewWorld(*debug)
	wo.Run(g.Run)
}
