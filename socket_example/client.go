package main

import (
	"bufio"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/:9020"
var url = "ws://localhost:9020/echo"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("map-id:prontera")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	scanner := bufio.NewScanner(ws)

	for scanner.Scan() {
		fmt.Printf("Receive: %s\n", scanner.Text())
	}

}
