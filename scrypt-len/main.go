package main

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

func PrefixScryptWithSalt(prefix, salt string, components ...[]byte) []byte {
	k, err := scrypt.Key([]byte(prefix), []byte(salt), 256, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	var res []byte
	res = append(res, k...)
	for _, c := range components {
		if len(c) == 0 {
			continue
		}
		hashed, err := scrypt.Key(c, []byte(salt), 256, 8, 1, 32)
		if err != nil {
			panic(err)
		}
		res = append(res, hashed...)
	}
	return res
}

func main() {
	components := [][]byte{
		[]byte("foo"), []byte("foo"), []byte("foo"), []byte("foo"), []byte("foo"), []byte("foo"),
		[]byte("foo"),
		[]byte("foo"),
		[]byte("foo"),
		[]byte("foo"),
		[]byte("foo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
		[]byte("fooooooooo"),
	}
	result := PrefixScryptWithSalt("prefix", "saltysalt", components...)
	fmt.Println("hash len", len(result))
}
