package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	playerStatusIdle   = iota
	playerStatusMoving // random movement
	playerStatusDead   // player has died
	playerStatusCombat // player is in combat
)

//Player struct
type Player struct {
	currentMap *Map
	mu         *sync.Mutex
	sockets    map[string]*websocket.Conn
	cmdChan    chan string
	statusChan chan string
	id         int64
	memID      string
	positionX  int64
	positionY  int64
	hp         int64
	maxHP      int64
	walkSpeed  int64
	status     int64
}

// Run the player
func (p *Player) Run() {
	go p.socket()
}

// ----------------- Combat ------------------ //

func (p *Player) damage(d int64) bool {
	if p.getStatus() == playerStatusDead {
		return false
	}

	hp := p.getHP()

	if hp > 0 {
		hp += d * -1
		p.setHP(hp)
		if hp <= 0 {
			go p.die()
			return false
		}

		p.cmdChan <- "p_dmg:" + strconv.FormatInt(d, 10)
		return true
	}

	return false
}

// ----------------- Die / Respawn ------------------ //

func (p *Player) die() {
	p.setStatus(playerStatusDead)
	p.cmdChan <- "p_die"
	<-time.NewTimer(time.Duration(random(2000, 3000)) * time.Millisecond).C
	p.respawn()
}

func (p *Player) respawn() {
	p.mu.Lock()
	p.hp = p.maxHP
	p.mu.Unlock()
	p.setXY(random(1, p.currentMap.size[0]), random(1, p.currentMap.size[1]))
	p.cmdChan <- "p_respawn"
	p.setStatus(playerStatusIdle)
}

// ----------------- getters / setters ------------------ //

func (p *Player) getHP() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.hp
}

func (p *Player) setHP(hp int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hp = hp
}

func (p *Player) setXY(x, y int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.positionX, p.positionY = x, y
}

func (p *Player) getXY() (int64, int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.positionX, p.positionY
}

func (p *Player) setStatus(status int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status = status
}

func (p *Player) getStatus() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.status
}

// ----------------- Basic Info ------------------ //

func (p *Player) getBasicInfo() []byte {
	dead := false
	if p.getStatus() == playerStatusDead {
		dead = true
	}
	p.mu.Lock()
	player := map[string]interface{}{
		"id":        p.id,
		"hp":        p.hp,
		"maxHP":     p.maxHP,
		"positionX": p.positionX,
		"positionY": p.positionY,
		"walkSpeed": p.walkSpeed,
		"dead":      dead,
	}
	p.mu.Unlock()
	b, _ := json.Marshal(player)
	return b
}

// ----------------- Sockets ------------------ //

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (p *Player) socket() {
	for {
		select {
		case cmd := <-p.cmdChan:
			jsonData := []byte{}

			switch cmd {
			case "p_move":
				x, y := p.getXY()
				jsonData, _ = json.Marshal([2]int64{x, y})
			case "p_die":
				jsonData, _ = json.Marshal(nil)
			case "p_respawn":
				jsonData, _ = json.Marshal(map[string]interface{}{
					"hp":        p.hp,
					"positionX": p.positionX,
					"positionY": p.positionY,
				})
			}

			if strings.HasPrefix(cmd, "p_dmg") {
				split := strings.Split(cmd, ":")
				cmd = split[0]
				jsonData, _ = json.Marshal(split[1])
			}

			go p.rangeSockets(cmd, string(jsonData))
		}
	}
}

func (p *Player) rangeSockets(cmd, data string) {
	data = fmt.Sprintf("%s:%s\n", cmd, data)
	p.mu.Lock()
	for _, sock := range p.sockets {
		go func(sock *websocket.Conn, data string) {
			sock.Write([]byte(data))
		}(sock, data)
	}
	p.mu.Unlock()
}

func (p *Player) addSocket(address string, ws *websocket.Conn) {
	p.mu.Lock()
	p.sockets[address] = ws
	p.mu.Unlock()
}

func (p *Player) delSocket(address string) {
	p.mu.Lock()
	delete(p.sockets, address)
	p.mu.Unlock()
}
