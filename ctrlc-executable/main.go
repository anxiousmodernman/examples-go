package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP, syscall.SIGQUIT)
	exit := longRunningGoroutine(quit)
	log.Println("Waiting for signal")
	code := <-exit
	log.Printf("Exited with code %v\n", code)
	os.Exit(code)

}

func longRunningGoroutine(quit chan os.Signal) chan int {
	exit := make(chan int)

	timer := time.NewTicker(2 * time.Second)

	go func() {

		for {
			select {
			case <-timer.C:
				fmt.Println("tick")
			case <-quit:
				fmt.Println("we out")
				exit <- 1
				break
			}
		}
	}()

	return exit
}
