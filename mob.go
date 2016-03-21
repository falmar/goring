package main

import (
	"fmt"
	"time"
)

// Monster basic interface
type Monster struct {
	id         int64
	name       string
	hp         int64
	level      int64
	sightRange int64
	walkRange  int64
	positionX  int64
	positionY  int64
	basemap    *Map
	walkSpeed  float64
}

var retry int

// NewMonster instance for maps
func NewMonster(id int64, basemap *Map) *Monster {
	return &Monster{
		id:         id,
		name:       "Poring",
		hp:         50,
		level:      1,
		sightRange: 12,
		walkRange:  6,
		basemap:    basemap,
		positionX:  basemap.size / 2,
		positionY:  basemap.size / 2,
		walkSpeed:  0.7,
	}
}

// Run start monster functions
func (m *Monster) Run() {
	go m.Move()
}

// Move the monster around the map
func (m *Monster) Move() {
	timer := time.NewTimer(time.Second * time.Duration(random(1, 5)))
	<-timer.C
	oX := m.positionX
	oY := m.positionY
	var nX int64
	var nY int64

	for {
		for {
			nX = oX + random(m.walkRange*-1, m.walkRange)
			if nX > 1 && nX < m.basemap.size {
				break
			}
		}

		for {
			nY = oY + random(m.walkRange*-1, m.walkRange)
			if nY > 1 && nY < m.basemap.size {
				break
			}
		}

		m.positionX = nX
		m.positionY = nY

		if m.routeWalk(oX, oY, nX, nY) {
			break
		}
	}

	fmt.Println(fmt.Sprintf("%s-%d moving to X:%d Y:%d", m.name, m.id, nX, nY))

	m.Move()
}

func (m *Monster) routeWalk(oX, oY, nX, nY int64) bool {
	var X, Y int64 = oX, oY
	var route [][]int64
	var movements int64

	for {

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
			return false
		}
		return m.routeWalk(oX, oY, nX, nY)
	}

	retry = 0
	return true
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
