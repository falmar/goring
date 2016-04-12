package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	mobStatusIdle = iota
	mobStatusMoving
	mobStatusCombat
)

// Monster basic interface
type Monster struct {
	mu         *sync.Mutex
	statusChan chan string
	cmdChan    chan string
	sockets    map[string]*websocket.Conn
	memID      string
	id         int64
	name       string
	hp         int64
	maxHP      int64
	level      int64
	baseMap    *Map
	positionX  int64
	positionY  int64
	sightRange int64
	walkRange  int64
	walkSpeed  float64
	walkRoute  []int64
	idleRange  []int64
	status     int64
	aggresive  bool
}

// NewMonster instance for maps
func NewMonster(id int64, baseMap *Map) *Monster {
	m := &Monster{
		mu:         &sync.Mutex{},
		statusChan: make(chan string),
		cmdChan:    make(chan string),
		sockets:    map[string]*websocket.Conn{},
		id:         id,
		name:       "Poring",
		hp:         50,
		level:      1,
		baseMap:    baseMap,
		positionX:  random(1, baseMap.size[0]),
		positionY:  random(1, baseMap.size[1]),
		walkRange:  6,
		walkSpeed:  1.6,
		idleRange:  []int64{2, 4},
		sightRange: 3,
		aggresive:  true,
		status:     mobStatusIdle,
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
					go m.sight()
				}
			}
		}
	}()

	go m.socket()
	go m.move()
	if m.aggresive {
		go m.sight()
	}

	loadedMob <- true
}

// ----------------- Random Movement ------------------ //

// Move the monster around the map
func (m *Monster) move() {
	idle := time.NewTimer(time.Duration(random(m.idleRange[0], m.idleRange[1])) * time.Second)
	<-idle.C
	oX, oY := m.getXY()
	var nX int64
	var nY int64

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

	for {
		if route, ok := m.routeMove(oX, oY, nX, nY); ok {
			walkTicker := time.NewTicker(time.Duration(m.walkSpeed*1000) * time.Millisecond)
			if m.getStatus() != mobStatusIdle {
				fmt.Println("Dont move...")
				return
			}
			m.setStatus(mobStatusMoving)
			for i := 0; i < len(route); i++ {
				if m.getStatus() != mobStatusMoving {
					fmt.Println("Stop moving")
					return
				}
				if i > 0 {
					<-walkTicker.C
				}
				m.mu.Lock()
				m.walkRoute = []int64{route[i][0], route[i][1]}
				m.mu.Unlock()
				m.setXY(route[i][0], route[i][1])
				m.cmdChan <- "move"
			}
			m.setStatus(mobStatusIdle)
			m.statusChan <- "idle"
			break
		}
	}
}

func (m Monster) routeMove(oX, oY, nX, nY int64) ([][]int64, bool) {
	var route [][]int64
	var movements int64
	var retry int

	for {
		route = [][]int64{}
		var X, Y int64 = oX, oY
		movements = 0

		for {
			X, Y = m.routeXY(X, nX, Y, nY)

			movements++
			if movements > m.walkRange {
				break
			}

			route = append(route, []int64{X, Y})

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

func (m Monster) routeXY(X, nX, Y, nY int64) (int64, int64) {
	r := random(1, 2)
	dr := random(1, 100)

	if (r == 1 || Y == nY) && X != nX {
		X = m.routeX(X, nX)
		if dr > 50 && dr < 100 {
			Y = m.routeX(Y, nY)
		}
	} else if (r == 2 || X == nX) && Y != nY {
		Y = m.routeY(Y, nY)
		if dr > 0 && dr < 50 {
			X = m.routeX(X, nX)
		}
	}

	return X, Y
}

func (m Monster) routeX(X, nX int64) int64 {
	if X > nX {
		X--
	} else if X < nX {
		X++
	}
	return X
}

func (m Monster) routeY(Y, nY int64) int64 {
	if Y > nY {
		Y--
	} else if Y < nY {
		Y++
	}
	return Y
}

// ----------------- sight ------------------ //

func (m *Monster) sight() {
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C

	if m.getStatus() == mobStatusCombat {
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

		var closestChan = make(chan int, len(targets))
		for i, t := range targets {
			tX, tY := t.getXY()
			go m.calculateClosestTarget(i, closestChan, x, y, tX, tY)
		}

		target := targets[<-closestChan]
		close(closestChan)

		fmt.Println("Target: ")
		fmt.Println(target)
	}

	//fmt.Println(targets)

	//route = m.calculateClosestRoute(x, y, nX, nY)

	if m.getStatus() != mobStatusCombat {
		m.statusChan <- "sight"
	}
}

// ----------------- Calculate closest Target & Route ------------------ //

func (m *Monster) calculateClosestTarget(i int, closestChan chan int, x, y, nX, nY int64) {
	for {
		x = m.routeX(x, nX)
		y = m.routeY(y, nY)

		if x == nX && y == nY {
			closestChan <- i
			break
		}
	}
}

func (m *Monster) calculateClosestRoute(oX, oY, tX, tY int64) [][2]int64 {
	var x, y int64 = oX, oY
	var nX, nY int64
	var attackRange int64 = 1
	var minX, minY, maxX, maxY int64 = oX - attackRange, oY - attackRange, oX + attackRange, oY + attackRange
	var route [][2]int64

	minMaxFunc := func(minX, maxX, minY, maxY, tX, tY int64) bool {
		if minX <= tX && maxX >= tX && minY <= tY && maxY >= tY {
			return true
		}
		return false
	}

	if !minMaxFunc(minX, maxX, minY, maxY, tX, tY) {
		if tX > oX {
			if (tX - oX) < attackRange {
				nX = oX + (tX - oX)
			} else {
				nX = tX - attackRange
			}
		} else {
			if (oX - tX) < attackRange {
				nX = oX - (oX - tX)
			} else {
				nX = tX + attackRange
			}
		}

		if tY > oY {
			if (tY - oY) < attackRange {
				nY = oY - (tY - oY)
			} else {
				nY = tY - attackRange
			}
		} else {
			if (oY - tY) < attackRange {
				nY = oY - (oY - tY)
			} else {
				nY = tY + attackRange
			}
		}

		minX, minY, maxX, maxY = nX-attackRange, nY-attackRange, nX+attackRange, nY+attackRange

		if minMaxFunc(minX, maxX, minY, maxY, nX, nY) {
			fmt.Println("Calculate route...")
		}

		for {
			x = m.routeX(x, nX)
			y = m.routeY(y, nY)
			route = append(route, [2]int64{x, y})
			if x == nX && y == nY {
				break
			}
		}
	}

	return route
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
	mob := map[string]interface{}{
		"id":        m.id,
		"hp":        m.hp,
		"positionX": m.positionX,
		"positionY": m.positionY,
		"walkSpeed": m.walkSpeed,
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
				route := m.walkRoute
				jsonData, _ = json.Marshal(route)
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
