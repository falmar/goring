package main

import "fmt"

var targets = [][2]int64{
	[2]int64{12, 8},
	[2]int64{13, 2},
	[2]int64{12, 4},
	[2]int64{11, 6},
	[2]int64{6, 5},
	[2]int64{9, 7},
}

func main() {
	var oX, oY int64 = 9, 9
	//var steps int
	var closestChan = make(chan map[int][][2]int64, 6)

	for i, t := range targets {
		go func(i int, closestChan chan map[int][][2]int64, x, y, nX, nY int64) {
			route := [][2]int64{}
			for {
				x = routeX(x, nX)
				y = routeY(y, nY)

				route = append(route, [2]int64{x, y})

				if x == nX && y == nY {
					closestChan <- map[int][][2]int64{
						i: route,
					}
					break
				}
			}
		}(i, closestChan, oX, oY, t[0], t[1])
	}

	resp := <-closestChan
	close(closestChan)

	for t, route := range resp {
		fmt.Println(t)
		fmt.Println(route)
	}

}

func calculateClosestRoute(oX, oY, tX, tY int64) [][2]int64 {
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
			x = routeX(x, nX)
			y = routeY(y, nY)

			route = append(route, [2]int64{x, y})

			if x == nX && y == nY {
				break
			}
		}
	}

	return route
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
