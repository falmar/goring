package main

import "encoding/json"

// Map generic struct
type Map struct {
	id   string
	name string
	size int64
	mobs []*Monster
}

// NewMap instance for server
func NewMap(id, name string, size int64) *Map {
	return &Map{
		id:   id,
		name: name,
		size: size,
	}
}

// Run map functionability
func (m *Map) Run() {
	m.loadMobs()
	// m.getStatus()
}

func (m *Map) loadMobs() {
	m.mobs = []*Monster{
		NewMonster(1002, m),
	}

	for _, mb := range m.mobs {
		go mb.Run()
	}
}

func (m *Map) getStatus() []byte {

	var mobs []int64
	for _, mb := range m.mobs {
		mobs = append(mobs, mb.id)
	}

	tmap := map[string]interface{}{
		"id":   m.id,
		"name": m.name,
		"size": m.size,
		"mobs": mobs,
	}
	b, _ := json.Marshal(tmap)

	return b
}
