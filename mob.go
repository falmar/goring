package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	mobStatusIdle         = iota
	mobStatusMoving       // random movement
	mobStatusMovingCombat // move to enter combat
	mobStatusMovingKin    // help another mob same kin
	mobStatusDead         // mob has died
	mobStatusCombat       // mob is in comabt
)

// Monster basic interface
type Monster struct {
	mu          *sync.Mutex
	statusChan  chan string
	cmdChan     chan string
	sockets     map[string]*websocket.Conn
	memID       string
	id          int64
	name        string
	hp          int64
	maxHP       int64
	level       int64
	baseMap     *Map
	positionX   int64
	positionY   int64
	sightRange  int64
	idleTime    [2]int64
	walkRange   int64
	walkSpeed   int64
	status      int64
	aggresive   bool
	respawnTime [2]int64
}

// NewMonster instance for maps
func NewMonster(id int64, baseMap *Map) *Monster {
	m := &Monster{
		mu:          &sync.Mutex{},
		statusChan:  make(chan string),
		cmdChan:     make(chan string),
		sockets:     map[string]*websocket.Conn{},
		id:          id,
		name:        "Poring",
		hp:          50,
		level:       1,
		baseMap:     baseMap,
		positionX:   random(1, baseMap.size[0]),
		positionY:   random(1, baseMap.size[1]),
		walkRange:   6,
		walkSpeed:   1600,
		idleTime:    [2]int64{2, 4},
		sightRange:  4,
		aggresive:   false,
		status:      mobStatusIdle,
		respawnTime: [2]int64{1, 10},
	}

	m.memID = fmt.Sprintf("%p", m)

	return m
}

// Run start monster functions
func (m *Monster) Run(loadedMob chan<- bool) {

	go func() {
		for {
			select {
			case cmd := <-m.statusChan:
				switch cmd {
				case "idle":
					go m.move()
				case "sight":
					if m.aggresive {
						go m.sight()
					}
				case "die":
					go m.die()
				case "respawn":
					go m.respawn()
				}
			}
		}
	}()

	go m.socket()
	m.spawnStatus()

	loadedMob <- true
}

// ----------------- Random Movement ------------------ //

// Move the monster around the map
func (m *Monster) move() {
	<-time.NewTimer(time.Duration(random(m.idleTime[0], m.idleTime[1])) * time.Second).C
	oX, oY := m.getXY()
	var nX int64
	var nY int64

	for {

		for {
			nX = oX + random(m.walkRange*-1, m.walkRange)
			if nX > 1 && nX < m.baseMap.size[0] {
				break
			}
		}

		for {
			nY = oY + random(m.walkRange*-1, m.walkRange)
			if nY > 1 && nY < m.baseMap.size[1] {
				break
			}
		}

		if m.getStatus() != mobStatusIdle || m.getStatus() == mobStatusDead {
			fmt.Println("Dont move...")
			return
		}

		if route, ok := m.routeRandomMove(oX, oY, nX, nY); ok {

			if m.walk(route, mobStatusMoving) {
				m.setStatus(mobStatusIdle)
				m.statusChan <- "idle"
			}

			break
		}
	}
}

func (m *Monster) walk(route [][]int64, requiredStatus int64) bool {
	m.setStatus(requiredStatus)
	var i int
	for i = 0; i < len(route); i++ {
		if m.getStatus() != requiredStatus {
			fmt.Println("Stop moving")
			return false
		}
		m.setXY(route[i][0], route[i][1])
		m.cmdChan <- "move"
		<-time.NewTimer(time.Duration(m.walkSpeed) * time.Millisecond).C
	}

	return len(route) == i
}

func (m Monster) routeRandomMove(oX, oY, nX, nY int64) ([][]int64, bool) {
	var route [][]int64
	var movements int64
	var retry int

	for {
		route = [][]int64{}
		var X, Y int64 = oX, oY
		movements = 0

		for {
			if movements > m.walkRange {
				break
			}

			X, Y = routeXY(X, nX, Y, nY)
			route = append(route, []int64{X, Y})
			movements++

			if X == nX && Y == nY {
				break
			}
		}

		if movements > m.walkRange {
			retry++
			if retry > 10 {
				return route, false
			}
			continue
		}

		break
	}

	return route, true
}

// ----------------- sight ------------------ //

func (m *Monster) sight() {
	<-time.NewTimer(100 * time.Millisecond).C

	currentStatus := m.getStatus()

	if currentStatus == mobStatusCombat || currentStatus == mobStatusDead {
		fmt.Println("Dont sight")
		return
	}

	x, y := m.getXY()

	var mapMaxX, mapMaxY int64 = m.baseMap.size[0], m.baseMap.size[1]
	var minX, minY, maxX, maxY int64

	if minX = x - m.sightRange; minX < 1 {
		minX = 1
	}
	if maxX = x + m.sightRange; maxX > mapMaxX {
		maxX = mapMaxX
	}
	if minY = y - m.sightRange; minY < 1 {
		minY = 1
	}
	if maxY = y + m.sightRange; maxY > mapMaxY {
		maxY = mapMaxY
	}

	players := m.baseMap.players
	var targets []*Player

	for _, p := range players {
		tX, tY := p.getXY()
		if minX <= tX && maxX >= tX && minY <= tY && maxY >= tY {
			targets = append(targets, p)
		}
	}

	if targets != nil {
		m.setStatus(mobStatusCombat)
		x, y = m.getXY()

		var closestTarget = make(chan *Player, len(targets))
		for _, t := range targets {
			go m.calculateClosestTarget(t, x, y, closestTarget)
		}

		//target :=
		<-closestTarget
		close(closestTarget)

		// Simulate combat
		//nX, nY := target.getXY()

		//route := m.calculateClosestRoute(x, y, nX, nY)

		return
	}

	if m.getStatus() != mobStatusCombat {
		m.statusChan <- "sight"
	}
}

// ----------------- Calculate closest target ------------------ //

func (m *Monster) calculateClosestTarget(t *Player, x, y int64, closestChan chan<- *Player) {
	nX, nY := t.getXY()
	for {
		x = routeX(x, nX)
		y = routeY(y, nY)

		if x == nX && y == nY {
			closestChan <- t
			break
		}
	}
}

// ----------------- Die / Respawn ------------------ //

func (m *Monster) die() {
	respawn := time.Duration(random(m.respawnTime[0], m.respawnTime[1])) * time.Second
	fmt.Println("Die.....; respawn in:", respawn)
	m.setStatus(mobStatusDead)
	m.cmdChan <- "die"
	<-time.NewTimer(respawn).C
	m.statusChan <- "respawn"
}

func (m *Monster) respawn() {
	fmt.Println("Respawning....")
	m.hp = m.maxHP
	m.positionX = random(1, m.baseMap.size[0])
	m.positionY = random(1, m.baseMap.size[1])
	m.setStatus(mobStatusIdle)
	m.spawnStatus()
}

func (m *Monster) spawnStatus() {
	m.statusChan <- "idle"
	m.statusChan <- "sight"
}

// ----------------- getters / setters ------------------ //

func (m *Monster) setXY(x, y int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.positionX, m.positionY = x, y
}

func (m *Monster) getXY() (int64, int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.positionX, m.positionY
}

func (m *Monster) setStatus(status int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = status
}

func (m *Monster) getStatus() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// ----------------- Basic Info ------------------ //

func (m *Monster) getBasicInfo() []byte {
	m.mu.Lock()
	dead := false
	if m.status == mobStatusDead {
		dead = true
	}
	mob := map[string]interface{}{
		"id":        m.id,
		"hp":        m.hp,
		"positionX": m.positionX,
		"positionY": m.positionY,
		"walkSpeed": m.walkSpeed,
		"dead":      dead,
	}
	m.mu.Unlock()
	b, _ := json.Marshal(mob)
	return b
}

// ----------------- Socket ------------------ //

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (m *Monster) socket() {
	for {
		select {
		case cmd := <-m.cmdChan:
			jsonData := []byte{}

			switch cmd {
			case "move":
				x, y := m.getXY()
				jsonData, _ = json.Marshal([2]int64{x, y})
			case "die":
				jsonData, _ = json.Marshal(nil)
			case "respawn":
				jsonData, _ = json.Marshal(map[string]interface{}{
					"hp":        m.hp,
					"positionX": m.positionX,
					"positionY": m.positionY,
				})
			}

			go m.rangeSockets(cmd, string(jsonData))
		}
	}
}

func (m *Monster) rangeSockets(cmd, data string) {
	data = fmt.Sprintf("%s:%s\n", cmd, data)
	m.mu.Lock()
	sockets := m.sockets
	m.mu.Unlock()
	for _, sock := range sockets {
		go func(sock *websocket.Conn, data string) {
			sock.Write([]byte(data))
		}(sock, data)
	}
}

func (m *Monster) addSocket(address string, ws *websocket.Conn) {
	m.mu.Lock()
	m.sockets[address] = ws
	m.mu.Unlock()
}

func (m *Monster) delSocket(address string) {
	m.mu.Lock()
	delete(m.sockets, address)
	m.mu.Unlock()
}
