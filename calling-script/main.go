package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("./printstuff.sh")

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("whoops %v", err)
	}
	fmt.Println(string(out))
}
