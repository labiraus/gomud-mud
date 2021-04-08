package world

import "github.com/labiraus/gomud-common/core"

type occupantRequest struct {
	room      *Room
	occupants chan<- map[uint]core.Entity
}

type move struct {
	origin      *Room
	destination uint
	ent         core.Entity
	dir         core.Direction
	done        chan<- bool
}

const (
	//SpawnPoint is where all characters start from
	SpawnPoint uint = 1
)
