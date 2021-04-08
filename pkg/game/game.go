package game

import (
	"context"
	"log"
	"strings"

	"./pkg/world"
	"github.com/labiraus/gomud-common/core"
)

//Game represents all logged in users
type Game struct {
	Users chan<- *User
	Ctx   context.Context 	

type listUpdate struct {
	u   *User
	add bool
}

//New creates the Game object
func New(ctx context.Context) (*Game, error) {
	w, err := world.New(ctx)
	if err != nil {
		return nil, err
	}
	userAddition := make(chan *User)
	updateList := make(chan listUpdate, 10)

	//Managed logged in user list
	go loggedInUsers(ctx, updateList)

	go handleNewUsers(ctx, userAddition, updateList, w)

	return &Game{
		Users: userAddition,
		Ctx:   ctx,
	}, nil
}

func handleNewUsers(ctx context.Context, userAddition <-chan *User, updateList chan<- listUpdate, w *world.World) {
	for {
		select {
		case <-ctx.Done():
			return
		case u := <-userAddition:
			go func() {
				u.start(func(u *User) {
					updateList <- listUpdate{u, true}
				})
				defer u.Cancel()
				go func() {
					<-u.Ctx.Done()
					updateList <- listUpdate{u, false}
				}()
				<-w.Spawn(u)

				for {
					text, err := u.read()
					if err != nil {
						log.Println(err)
						return
					}
					u.handleCommand(strings.ToLower(text[0]), text[1:])
					if err != nil {
						log.Println(err)
						return
					}
				}
			}()
		}
	}
}

func (u *User) handleCommand(command string, text []string) error {
	gAction, ok := globalActions[command]
	if ok {
		return gAction(u.Ctx, u, text)
	}
	dir, ok := core.Directions[command]
	if ok {
		<-u.location.Move(u.Ctx, u, dir)
		return nil
	}
	uAction, ok := u.actions[command]
	if ok {
		return uAction(u.Ctx, u, text)
	}

	return nil
}

func loggedInUsers(ctx context.Context, updateList <-chan listUpdate) {
	loggedInUsers := make(map[uint]*User)

	for {
		select {
		case <-ctx.Done():
			return
		case response := <-updateList:
			if response.u == nil || response.u.id == 0 {
				continue
			}
			if response.add {
				current, ok := loggedInUsers[response.u.id]
				if ok && current != nil {
					current.kill("Logged in a second time")
				}
				loggedInUsers[response.u.id] = response.u
			} else {
				delete(loggedInUsers, response.u.id)
			}
		}
	}
}
