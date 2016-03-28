package main

import "sync"

//Player struct
type Player struct {
	mu                   *sync.Mutex
	positionX, positionY int64
}

func (p *Player) getXY() (int64, int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.positionX, p.positionY
}
