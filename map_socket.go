package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func getMap(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}

	if i := strings.Index(string(msg[:n]), "map-id"); i == 0 {
		mapID := strings.Split(string(msg[:n]), ":")[1]
		if m, ok := maps[mapID]; ok {
			ws.Write(m.getBasicInfo())
			addr := fmt.Sprintf("%p", ws)
			m.addSocket(addr, ws)
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			defer m.delSocket(addr)
		}
	}
}
