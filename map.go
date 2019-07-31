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
	cmdChan chan string
	id      string
	name    string
	size    [2]int64
	mobs    map[string]*Monster
	sockets map[string]*websocket.Conn
	players map[string]*Player
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
		players: map[string]*Player{},
	}
}

// Run map functionability
func (m *Map) Run(loadedMap chan<- bool) {
	m.loadPlayers()
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
		mob.walkSpeed = 1200
		m.mobs[mob.memID] = mob
	}

	for i := 0; i < 2; i++ {
		mob := NewMonster(1007, m)
		mob.aggresive = true
		m.mobs[mob.memID] = mob
	}

	for i := 0; i < 2; i++ {
		mob := NewMonster(1049, m)
		mob.walkSpeed = 900
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

// fake func
func (m *Map) loadPlayers() {
	fmt.Println("Player at X:", 10, "Y:", 10)
	player := &Player{
		currentMap: m,
		mu:         &sync.Mutex{},
		positionX:  10,
		positionY:  10,
		walkSpeed:  1200,
		hp:         50,
		maxHP:      50,
		sockets:    map[string]*websocket.Conn{},
		cmdChan:    make(chan string),
		status:     playerStatusIdle,
	}

	player.memID = fmt.Sprintf("%p", player)

	m.players[player.memID] = player

	for _, p := range m.players {
		go p.Run()
	}
}

// ----------------- Basic Info ------------------ //

func (m *Map) getBasicInfo() []byte {
	var mobs []string
	for _, mb := range m.mobs {
		mobs = append(mobs, mb.memID)
	}

	var players []string
	for i := range m.players {
		players = append(players, i)
	}

	tmap := map[string]interface{}{
		"id":      m.id,
		"name":    m.name,
		"size":    m.size,
		"mobs":    mobs,
		"players": players,
	}

	b, _ := json.Marshal(tmap)

	return []byte(fmt.Sprintf("info:%s\n", string(b)))
}

func (m *Map) getMob(mobID string) (*Monster, bool) {
	mob, ok := m.mobs[mobID]
	return mob, ok
}

func (m *Map) getPlayer(playerID string) (*Player, bool) {
	player, ok := m.players[playerID]
	return player, ok
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
