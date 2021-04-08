package world

import (
	"fmt"

	"github.com/labiraus/gomud-common/core"
	"github.com/labiraus/gomud-common/utils"
)

type leave struct {
	ent  core.Entity
	dir  core.Direction
	room *Room
}

func (a leave) Act() <-chan struct{} {
	done := make(chan struct{})
	if a.room == nil {
		close(done)
		return done
	}

	go func() {
		defer close(done)
		select {
		case <-a.room.ctx.Done():
			return
		case occupants, ok := <-a.room.GetOccupants():
			if !ok {
				return
			}

			acc := make(chan (<-chan struct{}), len(occupants))
			for _, ent := range occupants {
				if ent != a.ent {
					acc <- ent.Send(fmt.Sprintf("%v %v", a.ent.LeaveString(), a.dir.EnterString()))
				}
			}

			close(acc)
			<-utils.DrainAccumulator(a.room.ctx, acc)
			return
		}
	}()

	return done
}
