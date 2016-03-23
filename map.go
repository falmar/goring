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
	mobs    map[string]*Monster
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
		mobs:    make(map[string]*Monster),
	}
}

// Run map functionability
func (m *Map) Run() {
	m.loadMobs()
	go m.Socket()
}

func (m *Map) loadMobs() {

	//TODO: Loop to load mobs from db or json file

	for i := 0; i < 3; i++ {
		mob := NewMonster(1002, m)
		m.mobs[mob.memID] = mob
	}

	for i := 0; i < 2; i++ {
		mob := NewMonster(1049, m)
		m.mobs[mob.memID] = mob
	}

	for _, mb := range m.mobs {
		go mb.Run()
	}
}

func (m *Map) getBasicInfo() []byte {
	var mobs []string
	for _, mb := range m.mobs {
		mobs = append(mobs, mb.memID)
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
	mob, ok := m.mobs[mobID]
	return mob, ok
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
