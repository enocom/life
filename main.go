package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func clearScreen() {
	print("\033[H\033[2J")
}

func main() {
	clearScreen()

	generation := &Generation{}
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
		time.Sleep(500 * time.Millisecond)
		clearScreen()
		generation.Reproduce()
	}

}
