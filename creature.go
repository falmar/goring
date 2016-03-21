package main

// Creature basic interface
type Creature interface {
	Walk(x, y uint8)
	Attack()
	Move()
	Help()
}
