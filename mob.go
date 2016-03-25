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
	walkRoute  [][]int64
	idleRange  []int64
	status     int64
	sockets    map[string]*websocket.Conn
	statusChan chan string
	cmdChan    chan string
}

// NewMonster instance for maps
func NewMonster(id int64, baseMap *Map) *Monster {
	m := &Monster{
		id:         id,
		name:       "Poring",
		hp:         50,
		level:      1,
		baseMap:    baseMap,
		positionX:  random(1, baseMap.size[0]),
		positionY:  random(1, baseMap.size[1]),
		walkRange:  6,
		walkSpeed:  1.6,
		idleRange:  []int64{2, 6},
		sightRange: 12,
		sockets:    map[string]*websocket.Conn{},
		statusChan: make(chan string),
		cmdChan:    make(chan string),
	}

	m.memID = fmt.Sprintf("%p", m)

	return m
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

// ----------------- Random Movement ------------------ //

// Move the monster around the map
func (m *Monster) move() {
	oX := m.positionX
	oY := m.positionY
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
			break
		}
	}
}

func (m *Monster) routeWalk(oX, oY, nX, nY int64) ([][]int64, bool) {
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

		//TODO: Fix route is never bigger than walkRange

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

func (m *Monster) routeXY(X, nX, Y, nY int64) (int64, int64) {
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

// ----------------- sight ------------------ //

func (m *Monster) sight() {

}

// ----------------- Basic Info ------------------ //

func (m *Monster) getBasicInfo() []byte {
	mob := map[string]interface{}{
		"id":        m.id,
		"hp":        m.hp,
		"positionX": m.positionX,
		"positionY": m.positionY,
		"walkSpeed": m.walkSpeed,
	}
	b, _ := json.Marshal(mob)
	return b
}

// ----------------- Socket ------------------ //

//Socket for map output data about the map commands and changes
// for JS browser to handle
func (m *Monster) Socket() {
	for {
		select {
		case cmd := <-m.cmdChan:

			jsonData := []byte{}

			switch cmd {
			case "move":
				jsonData, _ = json.Marshal(m.walkRoute)
			}

			m.rangeSockets(cmd, string(jsonData))
		}
	}
}

func (m *Monster) rangeSockets(cmd, data string) {
	data = fmt.Sprintf("%s:%s\n", cmd, data)
	for _, sock := range m.sockets {
		go func(sock *websocket.Conn, data string) {
			sock.Write([]byte(data))
		}(sock, data)
	}
}

func (m *Monster) addSocket(address string, ws *websocket.Conn) {
	m.sockets[address] = ws
}

func (m *Monster) delSocket(address string) {
	delete(m.sockets, address)
}
