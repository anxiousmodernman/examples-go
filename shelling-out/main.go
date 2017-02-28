package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	c := exec.Command("git", "status")
	r, err := c.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scnr := bufio.NewScanner(r)
	lines := make(chan string)
	done := make(chan bool)

	go func() {
		for scnr.Scan() {
			lines <- fmt.Sprintf("%s", scnr.Text())
		}
		close(lines)
		done <- true
	}()

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)
	go func() {
		errChan <- c.Wait()
	}()

	timeout := time.After(1 * time.Second)

	for {
		select {
		case <-done:
			fmt.Println("goodbye")
			os.Exit(0)
		case <-timeout:
			//return errors.New("command timed out")
			log.Fatal(err)
		case err := <-errChan:
			if err != nil {
				log.Fatal("errchanerr", err)
			}
			fmt.Println("closing errChan")
			errChan = nil
		case l := <-lines:
			fmt.Println("LINE:", l)
		}
	}

}
