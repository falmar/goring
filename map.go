package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

// Map generic struct
type Map struct {
	mu      *sync.Mutex
	id      string
	name    string
	size    [2]int64
	mobs    map[string]*Monster
	sockets map[string]*websocket.Conn
	cmdChan chan string
}

// NewMap instance for server
func NewMap(id, name string, size [2]int64) *Map {
	return &Map{
		mu:      &sync.Mutex{},
		id:      id,
		name:    name,
		size:    size,
		sockets: make(map[string]*websocket.Conn),
		mobs:    make(map[string]*Monster),
	}
}

// Run map functionability
func (m *Map) Run(loadedMap chan<- bool) {
	m.loadMobs()
	loadedMap <- true
	go m.Socket()
}

// ----------------- Monsters ------------------ //

func (m *Map) loadMobs() {
	var tMobs int
	var loadedMobs int
	loadedMob := make(chan bool)

	//TODO: Loop to load mobs from db or json file

	for i := 0; i < 3; i++ {
		mob := NewMonster(1002, m)
		m.mobs[mob.memID] = mob
	}

	for i := 0; i < 2; i++ {
		mob := NewMonster(1049, m)
		m.mobs[mob.memID] = mob
	}

	tMobs = len(m.mobs)

	for _, mb := range m.mobs {
		go mb.Run(loadedMob)
	}

mobCheckLoop:
	for {
		select {
		case <-loadedMob:
			loadedMobs++
			totalMobs++
			if loadedMobs == tMobs {
				break mobCheckLoop
			}
		}
	}
}

// ----------------- Basic Info ------------------ //

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

// ----------------- Socket Functions ------------------ //

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (m *Map) Socket() {
	for {
		select {
		case cmd := <-m.cmdChan:

			jsonData := []byte{}

			switch cmd {
			case "something":
			}

			m.rangeSockets(cmd, string(jsonData))
		}
	}
}

func (m *Map) rangeSockets(cmd, data string) {
	data = fmt.Sprintf("%s:%s\n", cmd, data)
	m.mu.Lock()
	for _, sock := range m.sockets {
		go func(sock *websocket.Conn, data string) {
			sock.Write([]byte(data))
		}(sock, data)
	}
	m.mu.Unlock()
}

func (m *Map) addSocket(address string, ws *websocket.Conn) {
	m.mu.Lock()
	m.sockets[address] = ws
	m.mu.Unlock()
}

func (m *Map) delSocket(address string) {
	m.mu.Lock()
	delete(m.sockets, address)
	m.mu.Unlock()
}
