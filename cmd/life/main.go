// +build !windows

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/enocom/life/pkg/life"
)

func main() {
	var c config
	flag.IntVar(&c.size, "size", 10, "the size of the game's dimensions")
	flag.DurationVar(&c.rate, "rate", time.Second, "the rate of generation refresh")
	flag.Parse()

	go listenForInterrupt()

	g := life.NewGame(
		life.WithBoardSize(c.size),
		life.WithGenerationRate(c.rate),
	)
	g.Start()
}

func listenForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Exiting...")
	os.Exit(0)
}

type config struct {
	size int
	rate time.Duration
}
