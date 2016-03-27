package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

func getMob(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Println(err)
	}

	if i := strings.Index(string(msg[:n]), "map-id"); i == 0 {
		mapID := strings.Split(string(msg[:n]), ":")[1]

		if m, ok := maps[mapID]; ok {
			n, err = ws.Read(msg)
			if err != nil {
				log.Println(err)
			}

			if i = strings.Index(string(msg[:n]), "mob-id"); i == 0 {
				mobID := strings.Split(string(msg[:n]), ":")[1]

				if mob, ok := m.getMob(mobID); ok {

					ws.Write([]byte(fmt.Sprintf("info:%s\n", mob.getBasicInfo())))

					addr := fmt.Sprintf("%p", ws)
					mob.addSocket(addr, ws)
					defer mob.delSocket(addr)

					scanner := bufio.NewScanner(ws)

					for scanner.Scan() {
						cmd := scanner.Text()
						switch cmd {
						case "cmd":

						}
					}
				}
			}
		}
	}
}
