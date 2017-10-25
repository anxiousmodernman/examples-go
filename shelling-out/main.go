package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	//c := exec.Command("bash", "-c", "git", "status")
	c := exec.Command("bash")
	r, err := c.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdin, _ := c.StdinPipe()

	scnr := bufio.NewScanner(r)
	lines := make(chan string)
	done := make(chan bool)

	go func() {
		fmt.Println("shelling out:  git status")
		for scnr.Scan() {
			lines <- fmt.Sprintf("%s", scnr.Text())
		}
		close(lines)
		done <- true
	}()

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	script := `
	ls
	echo foobar
	echo $PATH
	`
	io.Copy(stdin, bytes.NewBuffer([]byte(script)))

	timeout := time.After(1 * time.Second)

	for {
		select {
		case <-done:
			err := c.Wait()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("goodbye")
			os.Exit(0)
		case <-timeout:
			// return errors.New("command timed out")
			log.Fatal(err)
		case l := <-lines:
			fmt.Println("LINE:", l)
		}
	}

}
