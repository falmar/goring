package main

import "math/rand"

func random(min, max int64) int64 {
	return (rand.Int63n(max-min) + min) + 1
}

// Combat

func inAttackRange(x, y, tX, tY, attackRange int64) bool {
	var minX, minY, maxX, maxY int64

	minX = x - attackRange
	minY = y - attackRange
	maxX = x + attackRange
	maxY = y + attackRange

	if minX <= tX && maxX >= tX && minY <= tY && maxY >= tY {
		return true
	}

	return false
}

// Sight

func inSightRange(x, y, tX, tY, sightRange int64, m *Map) bool {
	minX, minY, maxX, maxY := minMaxSightRange(x, y, sightRange, m)

	if minX <= tX && maxX >= tX && minY <= tY && maxY >= tY {
		return true
	}

	return false
}

func minMaxSightRange(x, y, sightRange int64, m *Map) (minX, minY, maxX, maxY int64) {
	var mapMaxX, mapMaxY int64 = m.size[0], m.size[1]

	if minX = x - sightRange; minX < 1 {
		minX = 1
	}
	if maxX = x + sightRange; maxX > mapMaxX {
		maxX = mapMaxX
	}
	if minY = y - sightRange; minY < 1 {
		minY = 1
	}
	if maxY = y + sightRange; maxY > mapMaxY {
		maxY = mapMaxY
	}

	return minX, minY, maxX, maxY
}

// Routes

func calculateClosestRoute(oX, oY, tX, tY, attackRange int64) [][2]int64 {
	var x, y int64 = oX, oY
	var nX, nY int64
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
			for {
				x = routeX(x, nX)
				y = routeY(y, nY)
				route = append(route, [2]int64{x, y})
				if x == nX && y == nY {
					break
				}
			}
		}
	}

	return route
}

func routeXY(X, nX, Y, nY int64) (int64, int64) {
	r := random(1, 2)
	dr := random(1, 100)

	if (r == 1 || Y == nY) && X != nX {
		X = routeX(X, nX)
		if dr > 33 && dr < 100 {
			Y = routeX(Y, nY)
		}
	} else if (r == 2 || X == nX) && Y != nY {
		Y = routeY(Y, nY)
		if dr > 0 && dr < 67 {
			X = routeX(X, nX)
		}
	}

	return X, Y
}

func routeX(X, nX int64) int64 {
	if X > nX {
		X--
	} else if X < nX {
		X++
	}
	return X
}

func routeY(Y, nY int64) int64 {
	if Y > nY {
		Y--
	} else if Y < nY {
		Y++
	}
	return Y
}
