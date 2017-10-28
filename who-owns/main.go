package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// ls -l | awk '{print $3, $4 }'
	c := exec.Command("ls -l | awk '{print $3, $4}'")
	out, _ := c.Output()
	fmt.Println(string(out))

	fi, err := os.Stat("baz.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fi)
}
