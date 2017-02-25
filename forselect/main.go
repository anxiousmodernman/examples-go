package main

import (
	"errors"
	"fmt"
)

func main() {
	msgs := make(chan string)
	errs := make(chan error)

	go func() {
		msgs <- "hello"
		msgs <- "world"
		errs <- errors.New("nope")
	}()

	var running bool
	running = true

	for running {
		select {
		case m := <-msgs:
			fmt.Println("got msg", m)
		case err := <-errs:
			fmt.Println("got error", err)
			running = false
		}
	}
	fmt.Println("we are done")
}
