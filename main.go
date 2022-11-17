package main

import (
	"fmt"
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

	//go newConn(targetURL)
	go reuseConn(targetURL)

	var block chan bool
	<-block
}

func newConn(url string) {

	var transport *http.Transport

	for {
		t := time.Now().UTC()
		time.Sleep(1 * time.Second)
		transport = &http.Transport{}
		client := &http.Client{Transport: transport}

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(t, "newConn:", err)
			continue
		}

		fmt.Println(t, "newConn:", resp.Status)
		resp.Body.Close()
		transport.CloseIdleConnections()
	}

}

func reuseConn(url string) {

	client := &http.Client{}

	for {
		t := time.Now().UTC()
		time.Sleep(1 * time.Second)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(t, "oldConn:", err)
			continue
		}

		fmt.Println(t, "oldConn:", resp.Status)
		resp.Body.Close()
	}

}
