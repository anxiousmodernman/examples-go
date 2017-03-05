package main

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
)

func main() {
	iv := make([]byte, aes.BlockSize)

	fmt.Println(iv)
	rand.Read(iv)
	fmt.Println(iv)

}
