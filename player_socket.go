package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

func getPlayer(ws *websocket.Conn) {
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

			if i = strings.Index(string(msg[:n]), "player-id"); i == 0 {
				playerID := strings.Split(string(msg[:n]), ":")[1]

				if player, ok := m.getPlayer(playerID); ok {

					ws.Write([]byte(fmt.Sprintf("info:%s\n", player.getBasicInfo())))

					addr := fmt.Sprintf("%p", ws)
					player.addSocket(addr, ws)
					defer player.delSocket(addr)

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
