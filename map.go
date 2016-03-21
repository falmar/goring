package main

// Map generic struct
type Map struct {
	name string
	size int64
	mobs []*Monster
}

// NewMap instance for server
func NewMap(name string, size int64) *Map {
	return &Map{
		name: name,
		size: size,
	}
}

// Run map functionability
func (m *Map) Run() {
	m.loadMobs()
}

func (m *Map) loadMobs() {
	m.mobs = []*Monster{
		NewMonster(1002, m),
	}

	for _, mb := range m.mobs {
		go mb.Run()
	}
}
