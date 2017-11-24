package main

import (
	"github.com/explodes/go-wo"
	"github.com/explodes/go-wo/examples/flappy/internal"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	g := internal.NewFlappyWorld()
	wo.Run(g.Run)
}
