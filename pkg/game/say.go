package game

import (
	"context"
	"fmt"
	"strings"

	"github.com/labiraus/gomud-common/core"
	"github.com/labiraus/gomud-common/utils"
)

type sayAction struct {
	ctx  context.Context
	text []string
	ent  core.Entity
}

func (s sayAction) Act() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		select {
		case <-s.ctx.Done():
			return
		case occupants, ok := <-s.ent.Location().GetOccupants():
			if !ok {
				return
			}
			eName := s.ent.Name()

			acc := make(chan (<-chan struct{}), len(occupants))
			if len(s.text) > 2 && s.text[0] == "to" {
				target := utils.Search(occupants, s.text[1])
				if target != nil {
					text := strings.Join(s.text[2:], " ")
					tName := target.Name()

					for _, ent := range occupants {
						switch {
						case ent == s.ent:
							acc <- ent.Send(fmt.Sprintf("You say to %v: %v", tName, text))
						case ent == target:
							acc <- ent.Send(fmt.Sprintf("%v says to you: %v", eName, text))
						default:
							acc <- ent.Send(fmt.Sprintf("%v says to %v : %v", eName, tName, text))
						}
					}
					close(acc)
					<-utils.DrainAccumulator(s.ctx, acc)
					return
				}
			}
			text := strings.Join(s.text, " ")

			for _, ent := range occupants {
				switch {
				case ent == s.ent:
					acc <- ent.Send(fmt.Sprintf("You say: %v", text))
				default:
					acc <- ent.Send(fmt.Sprintf("%v says: %v", eName, text))
				}
			}
			close(acc)
			<-utils.DrainAccumulator(s.ctx, acc)
			return

		}
	}()

	return done
}

func say(ctx context.Context, ent core.Entity, text []string) error {
	if len(text) == 0 {
		return nil
	}
	ent.Location().Act(ctx, sayAction{
		ctx:  ctx,
		text: text,
		ent:  ent,
	})
	return nil
}
