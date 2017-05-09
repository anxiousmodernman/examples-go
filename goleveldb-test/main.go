package main

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		log.Fatal(err)
	}
	k := []byte("key")
	v := []byte("value")
	db.Put(k, v, nil)
	got, err := db.Get(k, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get:", got)
	fmt.Println("Get (string):", string(got))
}
