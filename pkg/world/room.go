package world

import (
	"context"
	"strings"

	"github.com/labiraus/gomud-common/core"
)

//Room is the location that an entity occupies
type Room struct {
	id          uint
	name        string
	description string
	ctx         context.Context
	roomActions chan<- core.Action
	neighbours  map[core.Direction]*Room
	world       *World
}

//ID returns a room's ID
func (r Room) ID() uint {
	return r.id
}

//Name returns a room's name
func (r Room) Name() string {
	return r.name
}

//Description returns a room's name
func (r Room) Description() string {
	return r.description
}

//GetOccupants gets a map of all occupants in a room
func (r *Room) GetOccupants() <-chan map[uint]core.Entity {
	out := make(chan map[uint]core.Entity)
	select {
	case <-r.ctx.Done():
	default:
		r.world.occupants <- occupantRequest{room: r, occupants: out}
	}
	return out
}

//Move moves an entity in a direction, returns false if there is no exit in that direction
func (r *Room) Move(ctx context.Context, ent core.Entity, d core.Direction) <-chan bool {
	out := make(chan bool, 1)

	go func() {
		destination, ok := r.neighbours[d]
		if ok {
			select {
			case <-r.ctx.Done():
			case <-ctx.Done():
			case r.world.moves <- move{
				origin:      r,
				destination: destination.id,
				dir:         d,
				ent:         ent,
				done:        out,
			}:
			}
		} else {
			select {
			case <-r.ctx.Done():
			case <-ctx.Done():
			case out <- false:
			}
			close(out)
		}
	}()

	return out
}

//Act pushes an action onto the room's action queue
func (r *Room) Act(ctx context.Context, a core.Action) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case r.roomActions <- a:
		case <-ctx.Done():
		case <-r.ctx.Done():
		}
	}()
	return done
}

//Exits lists all the exits to a room
func (r *Room) Exits() string {
	if len(r.neighbours) == 0 {
		return "No Exits"
	}
	if len(r.neighbours) == 1 {
		for exitID := range r.neighbours {
			return "Exit leading " + exitID.String()
		}
	}

	var exitArray [core.DirectionCount]string
	for exitID := range r.neighbours {
		exitArray[uint(exitID)] = exitID.String()
	}
	exitSlice := make([]string, 0, len(r.neighbours))
	for _, exit := range exitArray {
		if exit != "" {
			exitSlice = append(exitSlice, exit)
		}
	}
	if len(exitSlice) == 2 {
		return "Exits leading " + strings.Join(exitSlice, " and ")
	}
	return "Exits leading " + strings.Join(exitSlice[:len(exitSlice)-1], ", ") + ", and " + exitSlice[len(exitSlice)-1]
}
