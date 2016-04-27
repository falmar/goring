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
	mobStatusCombat       // mob is in combat
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
	attackPower [2]int64
	attackRange int64
	attackSpeed int64
	sightRange  int64
	sightTime   int64
	idleTime    [2]int64
	walkRange   int64
	walkSpeed   int64
	status      int64
	aggresive   bool
	respawnTime [2]int64
	target      *Player
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
		idleTime:    [2]int64{1000, 2600},
		sightRange:  6,
		sightTime:   100,
		aggresive:   true,
		status:      mobStatusIdle,
		respawnTime: [2]int64{1, 10},
		attackPower: [2]int64{7, 10},
		attackSpeed: 1870,
		attackRange: 1,
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
	<-time.NewTimer(time.Duration(random(m.idleTime[0], m.idleTime[1])) * time.Millisecond).C
	var oX, oY int64 = m.getXY()
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

func (m *Monster) walk(route [][2]int64, requiredStatus int64) bool {
	m.setStatus(requiredStatus)
	var i int
	for i = 0; i < len(route); i++ {
		if i < len(route) {
			<-time.NewTimer(time.Duration(m.walkSpeed) * time.Millisecond).C
		}
		if m.getStatus() != requiredStatus {
			return false
		}
		m.setXY(route[i][0], route[i][1])
		m.cmdChan <- "move"
	}

	return len(route) == i
}

func (m Monster) routeRandomMove(oX, oY, nX, nY int64) ([][2]int64, bool) {
	var route [][2]int64
	var movements int64
	var retry int

	for {
		route = [][2]int64{}
		var X, Y int64 = oX, oY
		movements = 0

		for {
			if movements > m.walkRange {
				break
			}

			X, Y = routeXY(X, nX, Y, nY)
			route = append(route, [2]int64{X, Y})
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
	<-time.NewTimer(time.Duration(m.sightTime) * time.Millisecond).C

	if m.getStatus() == mobStatusDead {
		return
	}

	x, y := m.getXY()

	if m.getTarget() == nil {

		var targets []*Player
		for _, p := range m.baseMap.players {
			if p.getStatus() == playerStatusDead {
				continue
			}
			tX, tY := p.getXY()
			if inSightRange(x, y, tX, tY, m.sightRange, m.baseMap) {
				targets = append(targets, p)
			}
		}

		if len(targets) > 0 {
			//TODO: Fix status
			m.setStatus(mobStatusCombat)
			x, y = m.getXY()

			var closestTarget = make(chan *Player, len(targets))
			for _, t := range targets {
				go m.calculateClosestTarget(t, x, y, closestTarget)
			}

			m.setTarget(<-closestTarget)
			close(closestTarget)

			go m.attack()

			return
		}

		m.statusChan <- "sight"
	}
}

// ----------------- Combat ------------------ //

func (m *Monster) attack() {

	for {
		target := m.getTarget()

		if target == nil || target.getStatus() == playerStatusDead {
			m.setTarget(nil)
			m.setStatus(mobStatusIdle)
			go m.spawnStatus()
			return
		}

		x, y := m.getXY()
		tX, tY := target.getXY()
		if inAttackRange(x, y, tX, tY, m.attackRange) {
			<-time.NewTimer(time.Duration(m.attackSpeed) * time.Millisecond).C
			tX, tY = target.getXY()
			if inAttackRange(x, y, tX, tY, m.attackRange) {
				if !target.damage(random(m.attackPower[0], m.attackPower[1])) {
					m.setTarget(nil)
				}
			}
			continue
		} else if inSightRange(x, y, tX, tY, m.sightRange, m.baseMap) {
			route := calculateClosestRoute(x, y, tX, tY, m.attackRange)
			m.walk(route, mobStatusMovingCombat)
			continue
		}

		m.setTarget(nil)
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
	m.setStatus(mobStatusDead)
	m.cmdChan <- "die"
	<-time.NewTimer(respawn).C
}

func (m *Monster) respawn() {
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

func (m *Monster) setTarget(t *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.target = t
}

func (m *Monster) getTarget() *Player {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.target
}

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
	dead := false
	if m.getStatus() == mobStatusDead {
		dead = true
	}
	m.mu.Lock()
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
