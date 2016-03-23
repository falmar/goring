package main

import (
	"encoding/json"
	"fmt"
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
	walkRoute  [][]int64
	idleRange  []int64
	status     int64
	sockets    map[string]*websocket.Conn
	statusChan chan string
	cmdChan    chan string
}

var retry int

// NewMonster instance for maps
func NewMonster(id int64, baseMap *Map) *Monster {
	return &Monster{
		id:         id,
		name:       "Poring",
		hp:         50,
		level:      1,
		baseMap:    baseMap,
		positionX:  random(1, baseMap.size),
		positionY:  random(1, baseMap.size),
		walkRange:  6,
		walkSpeed:  1.2,
		idleRange:  []int64{1, 4},
		sightRange: 12,
		sockets:    map[string]*websocket.Conn{},
		statusChan: make(chan string),
		cmdChan:    make(chan string),
	}
}

// Run start monster functions
func (m *Monster) Run() {

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

	go m.move()
	go m.sight()
	go m.Socket()
}

// Move the monster around the map
func (m *Monster) move() {
	oX := m.positionX
	oY := m.positionY
	var nX int64
	var nY int64

	for {
		nX = oX + random(m.walkRange*-1, m.walkRange)
		if nX > 1 && nX < m.baseMap.size {
			break
		}
	}

	for {
		nY = oY + random(m.walkRange*-1, m.walkRange)
		if nY > 1 && nY < m.baseMap.size {
			break
		}
	}

	if route, ok := m.routeWalk(oX, oY, nX, nY); ok {
		m.cmdChan <- "move"
		m.walkRoute = route
		walkTicker := time.NewTicker(time.Duration(m.walkSpeed*1000) * time.Millisecond)
		m.status = mobStatusMoving
		for i := 0; i < len(route); i++ {
			if m.status != mobStatusMoving {
				return
			}
			<-walkTicker.C
			m.positionX = route[i][0]
			m.positionY = route[i][1]
		}
		m.status = mobStatusIdle
		m.statusChan <- "idle"
	}
}

func (m *Monster) sight() {

}

func (m *Monster) routeWalk(oX, oY, nX, nY int64) ([][]int64, bool) {
	var X, Y int64 = oX, oY
	var route [][]int64
	var movements int64

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

	if int64(len(route)) > m.walkRange {
		retry++
		if retry > 6 {
			return route, false
		}
		return m.routeWalk(oX, oY, nX, nY)
	}

	retry = 0
	return route, true
}

func (m *Monster) routeXY(X, nX, Y, nY int64) (int64, int64) {
	r := random(1, 2)
	dr := random(1, 100)

	if (r == 1 || Y == nY) && X != nX {
		X = m.routeX(X, nX)
		if dr > 65 && dr < 100 {
			Y = m.routeX(Y, nY)
		}
	} else if (r == 2 || X == nX) && Y != nY {
		Y = m.routeY(Y, nY)
		if dr > 0 && dr < 45 {
			X = m.routeX(X, nX)
		}
	}

	return X, Y
}

func (m *Monster) routeX(X, nX int64) int64 {
	if X > nX {
		X--
	} else if X < nX {
		X++
	}
	return X
}

func (m *Monster) routeY(Y, nY int64) int64 {
	if Y > nY {
		Y--
	} else if Y < nY {
		Y++
	}
	return Y
}

func (m *Monster) getBasicInfo() []byte {
	mob := map[string]interface{}{
		"hp":        m.hp,
		"positionX": m.positionX,
		"positionY": m.positionY,
	}
	b, _ := json.Marshal(mob)
	return b
}

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (m *Monster) Socket() {
	for {
		select {
		case cmd := <-m.cmdChan:
			switch cmd {
			case "move":
				//fmt.Println(fmt.Sprintf("%s (%p) Moving to %d %d", m.name, m, m.positionX, m.positionY))
				data := map[string]interface{}{
					"route": m.walkRoute,
				}
				json, _ := json.Marshal(data)
				m.rangeSockets(cmd, string(json))
			}
		}
	}
}

func (m *Monster) rangeSockets(cmd, data string) {
	for _, sock := range m.sockets {
		sock.Write([]byte(fmt.Sprintln(cmd, ":", data)))
	}
}

func (m *Monster) addSocket(address string, ws *websocket.Conn) {
	m.sockets[address] = ws
}

func (m *Monster) delSocket(address string) {
	delete(m.sockets, address)
}
