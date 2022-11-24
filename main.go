package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: conn-reuse <URL>")
		os.Exit(1)
	}

	targetURL := os.Args[1]

	go newConn(targetURL)
	go reuseConn(targetURL)

	var block chan bool
	<-block
}

func req(client *http.Client, url string, prefix string, errors *int) {
	resp, err := client.Get(url)
	t := time.Now().UTC()
	if err != nil {
		*errors = *errors + 1
		fmt.Println(t, prefix, err, *errors)
		return
	}
	text, _ := io.ReadAll(resp.Body)
	fmt.Println(t, prefix, resp.Status, string(text))
	resp.Body.Close()
}

func newConn(url string) {

	newConnKeyLog, err := os.Create("newConn.keylog")

	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		KeyLogWriter: newConnKeyLog,
	}

	var transport *http.Transport
	errors := 0

	for {
		time.Sleep(1 * time.Second)
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client := &http.Client{Transport: transport}
		req(client, url, "newConn", &errors)
		transport.CloseIdleConnections()
	}

}

func reuseConn(url string) {

	oldConnKeyLog, err := os.Create("oldConn.keylog")

	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		KeyLogWriter: oldConnKeyLog,
	}

	transport := &http.Transport{
		MaxConnsPerHost: 1,
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}
	errors := 0

	for {
		time.Sleep(1 * time.Second)
		req(client, url, "oldConn", &errors)
	}

}
