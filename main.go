package main

import (
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

func req(client *http.Client, url string, prefix string) {
	resp, err := client.Get(url)
	t := time.Now().UTC()
	if err != nil {
		fmt.Println(t, prefix, err)
		return
	}
	text, _ := io.ReadAll(resp.Body)
	fmt.Println(t, prefix, resp.Status, string(text))
	resp.Body.Close()
}

func newConn(url string) {
	var transport *http.Transport

	for {
		time.Sleep(1 * time.Second)
		transport = &http.Transport{}
		client := &http.Client{Transport: transport}
		req(client, url, "newConn")
		transport.CloseIdleConnections()
	}

}

func reuseConn(url string) {

	transport := &http.Transport{
		MaxConnsPerHost: 1,
	}
	client := &http.Client{Transport: transport}

	for {
		time.Sleep(1 * time.Second)
		req(client, url, "oldConn")
	}

}
