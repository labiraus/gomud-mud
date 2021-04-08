package world

import (
	"context"
	"fmt"

	"github.com/labiraus/gomud-common/core"
	"github.com/labiraus/gomud-common/db"
)

//World contains all of the Rooms
type World struct {
	occupants chan<- occupantRequest
	moves     chan<- move
}

//New initialise the World
func New(ctx context.Context) (*World, error) {
	occupants := make(map[uint]map[uint]core.Entity)
	moves := make(chan move, 10)
	occupantRequests := make(chan occupantRequest, 10)
	w := &World{
		moves:     moves,
		occupants: occupantRequests,
	}

	fmt.Println("Building world")
	roomMap := make(map[uint]*Room)
	getRoom := func(id uint) *Room {
		room, ok := roomMap[id]
		if !ok {

			room = &Room{
				id:          id,
				roomActions: roomActions(ctx),
				ctx:         ctx,
				neighbours:  make(map[core.Direction]*Room),
				world:       w,
			}
			roomMap[id] = room

			occupants[id] = make(map[uint]core.Entity)
		}
		return room
	}

	roomSlice := db.GetRooms()
	for _, room := range roomSlice {
		newRoom := getRoom(room.ID)
		newRoom.name = room.Name
		newRoom.description = room.Description
		for _, connection := range room.Connections {
			neighbour := getRoom(connection.NeighbourID)
			newRoom.neighbours[connection.Dir] = neighbour
			if connection.TwoWay {
				neighbour.neighbours[connection.Dir.Reverse()] = newRoom
			}
		}
	}

	fmt.Println("World built")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case request := <-occupantRequests:
				request.occupants <- occupants[request.room.id]

				close(request.occupants)

			case move := <-moves:
				if move.origin != nil {
					delete(occupants[move.origin.id], move.ent.ID())
					move.origin.roomActions <- leave{ent: move.ent, dir: move.dir, room: move.origin}
				}

				if move.destination != 0 {
					destination, ok := roomMap[move.destination]
					room, ok := occupants[move.destination]
					if ok {
						room[move.ent.ID()] = move.ent
						destination.roomActions <- enter{ent: move.ent, dir: move.dir.Reverse(), room: destination}
					}
					move.ent.SetLocation(destination)
					move.done <- ok
				}

				close(move.done)
			}
		}
	}()

	return w, nil
}

//Spawn spawns an entity according to its Home()
func (w *World) Spawn(ent core.Entity) <-chan bool {
	done := make(chan bool, 1)
	w.moves <- move{ent: ent, destination: ent.Home(), done: done}
	return done
}

//roomActions creates an action channel
func roomActions(ctx context.Context) chan<- core.Action {
	actionChan := make(chan core.Action, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case action := <-actionChan:
				action.Act()
			}
		}
	}()
	return actionChan
}
