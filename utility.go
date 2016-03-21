package main

import "math/rand"

func random(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
