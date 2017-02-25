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

	for {
		select {
		case m := <-msgs:
			fmt.Println("got msg", m)
		case err := <-errs:
			fmt.Println("got error", err)
			goto Finished
		}
	}
Finished:
	fmt.Println("we are done")
}
