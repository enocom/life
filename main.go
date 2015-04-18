package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func clearScreen() {
	print("\033[H\033[2J")
}

func main() {
	var generationTime *int = flag.Int("gtime", 1000, "time in milliseconds between generations")
	var generationHeight *int = flag.Int("height", 20, "height of the generation")
	var generationWidth *int = flag.Int("width", 20, "height of the generation")
	flag.Parse()

	clearScreen()

	generation := &Generation{height: *generationHeight, width: *generationWidth}
	generation.Awaken()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()

	for i := 0; true; i++ {
		fmt.Println("Generation", i)
		fmt.Println(generation.ToString())
		time.Sleep(time.Duration(*generationTime) * time.Millisecond)
		clearScreen()
		generation.Reproduce()
	}

}
