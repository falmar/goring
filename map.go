package main

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

// Map generic struct
type Map struct {
	id      string
	name    string
	size    int64
	mobs    []*Monster
	sockets map[string]*websocket.Conn
}

// NewMap instance for server
func NewMap(id, name string, size int64) *Map {
	return &Map{
		id:      id,
		name:    name,
		size:    size,
		sockets: make(map[string]*websocket.Conn),
	}
}

// Run map functionability
func (m *Map) Run() {
	m.loadMobs()
	go m.fakeCmd()
	go m.Socket()
	go m.checkSockets()
}

func (m *Map) checkSockets() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for {
			<-ticker.C
			fmt.Println(m.sockets)
		}
	}()
}

func (m *Map) loadMobs() {
	m.mobs = []*Monster{
		NewMonster(1002, m),
	}

	for _, mb := range m.mobs {
		go mb.Run()
	}
}

func (m *Map) getBasicInfo() []byte {
	var mobs []int64
	for _, mb := range m.mobs {
		mobs = append(mobs, mb.id)
	}

	tmap := map[string]interface{}{
		"id":   m.id,
		"name": m.name,
		"size": m.size,
		"mobs": mobs,
	}

	b, _ := json.Marshal(tmap)

	return []byte(fmt.Sprintln(b))
}

//Socket for map
func (m *Map) Socket() {
	for {
		select {
		case cmd := <-mapCmdChan[m.id]:
			for _, sock := range m.sockets {
				sock.Write([]byte(fmt.Sprintln(cmd)))
			}
		}
	}
}

func (m *Map) addSocket(address string, ws *websocket.Conn) {
	m.sockets[address] = ws
}

func (m *Map) delSocket(address string) {
	delete(m.sockets, address)
}

func (m *Map) fakeCmd() {

	for {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C

		mapCmdChan[m.id] <- "Pix logged in"

		timer = time.NewTimer(3 * time.Second)
		<-timer.C

		mapCmdChan[m.id] <- "Pix logged out"
	}

}
