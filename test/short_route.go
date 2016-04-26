package main

import "fmt"

func main() {
	var x, y int64 = 9, 9
	var oX, oY, tX, tY int64 = 9, 9, 4, 18
	var nX, nY int64
	var attackRange int64 = 1
	var minX, minY, maxX, maxY int64 = oX - attackRange, oY - attackRange, oX + attackRange, oY + attackRange
	var route [][]int64

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
			y = routeX(y, nY)

			route = append(route, []int64{x, y})

			if x == nX && y == nY {
				break
			}
		}

	}

	fmt.Println(route)

	fmt.Println(tX, tY, nX, nY)
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
