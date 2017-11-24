package main

import (
	"flag"

	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/soccer/internal"
	"github.com/sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "enable debug drawing")
)

func main() {
	flag.Parse()
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	g := internal.NewWorld(*debug)
	wo.Run(g.Run)
}
