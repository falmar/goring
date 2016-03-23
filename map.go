package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

// Map generic struct
type Map struct {
	id      string
	name    string
	size    int64
	mobs    []*Monster
	sockets map[string]*websocket.Conn
	cmdChan chan string
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
	go m.Socket()
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

	return []byte(fmt.Sprintf("info:%s\n", string(b)))
}

func (m *Map) getMob(mobID string) (*Monster, bool) {
	return nil, false
}

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (m *Map) Socket() {
	for {
		select {
		case cmd := <-m.cmdChan:
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
