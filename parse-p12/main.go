package main

import (
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/pkcs12"
)

var (
	filename = flag.String("f", "", "pkcs12 file to read")
)

func main() {
	flag.Parse()
	if *filename == "" {
		log.Fatal("you must pass the pkcs12 file with -f flag")
	}
	fmt.Println("reading", *filename)
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}

	blocks, err := pkcs12.ToPEM(data, "password")
	if err != nil {
		panic(err)
	}

	var pemData []byte
	log.Printf("decoded %v blocks\n", len(blocks))
	for i, b := range blocks {
		log.Printf("block %v type %s\n", i, b.Type)
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	// then use PEM data for tls to construct tls certificate:
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	_ = config

}
