package main

import (
	"fmt"
	"sync"
)

//Player struct
type Player struct {
	mu        *sync.Mutex
	positionX int64
	positionY int64
	hp        int64
	dead      bool
}

// ----------------- Combat ------------------ //

func (p *Player) damage(d int64) {

	if p.dead {
		return
	}

	hp := p.getHP()

	if hp > 0 {
		hp += d * -1

		if hp <= 0 {
			fmt.Printf("Dead player %p\n", p)
			hp = 0
			p.dead = true
		}

		fmt.Println("DMG:", d, hp)

		p.setHP(hp)
	}
}

// ----------------- getters / setters ------------------ //

func (p *Player) getHP() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.hp
}

func (p *Player) setHP(hp int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hp = hp
}

func (p *Player) getXY() (int64, int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.positionX, p.positionY
}
