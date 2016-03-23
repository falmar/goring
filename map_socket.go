package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

func getMap(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Println(err)
	}

	if i := strings.Index(string(msg[:n]), "map-id"); i == 0 {
		mapID := strings.Split(string(msg[:n]), ":")[1]
		if m, ok := maps[mapID]; ok {

			ws.Write(m.getBasicInfo())

			addr := fmt.Sprintf("%p", ws)
			m.addSocket(addr, ws)
			defer m.delSocket(addr)

			scanner := bufio.NewScanner(ws)

			for scanner.Scan() {
				cmd := scanner.Text()
				switch cmd {
				case "some-cmd":
				}
			}
		}
	}
}
