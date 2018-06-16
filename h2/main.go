package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

func main() {
	//url := "https://www.microsoft.com/en-gb/"
	url := "https://http2.golang.org"

	clientProtocols := []string{"h2"}
	clientTlsConfig := &tls.Config{
		NextProtos: clientProtocols,
		//CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	}

	client := http.Client{
		Transport: &http2.Transport{TLSClientConfig: clientTlsConfig},
	}

	resp, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get protocol version
	fmt.Println("proto:", resp.Proto)
	fmt.Println("other proto:", resp.Request.Proto)

}
