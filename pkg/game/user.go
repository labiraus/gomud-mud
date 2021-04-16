package game

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/labiraus/gomud-common/core"
	"github.com/labiraus/gomud-mud/pkg/world"
)

//User is a single person logged in
type User struct {
	id          uint
	RecieveChan <-chan string
	SendChan    chan<- string
	Ctx         context.Context
	Cancel      context.CancelFunc
	name        string
	location    core.Location
	actions     map[string]func(ctx context.Context, ent core.Entity, text []string) error
}

//Init sets up a user
func (u *User) Init(ctx context.Context) core.Entity {
	u.Ctx, u.Cancel = context.WithCancel(ctx)
	u.actions = map[string]func(ctx context.Context, ent core.Entity, text []string) error{
		"say": say,
	}
	return u
}

//ID returns the User's ID
func (u *User) ID() uint {
	return u.id
}

//Name returns the User's character name
func (u *User) Name() string {
	return u.name
}

//Home returns the User's last known location
func (u *User) Home() uint {
	if u.location != nil {
		return u.location.ID()
	}
	return world.SpawnPoint
}

//EnterString returns a string when a User enters a room
func (u *User) EnterString() string {
	return fmt.Sprintf("%v enters", u.name)
}

//LeaveString returns a string when a User leaves a room
func (u *User) LeaveString() string {
	return fmt.Sprintf("%v leaves", u.name)
}

//SetLocation sets the User's location
func (u *User) SetLocation(location core.Location) {
	u.location = location
	u.Send("Entering " + location.Name())
	u.Send(location.Exits())
}

//Location returns a user's current location
func (u *User) Location() core.Location {
	return u.location
}

func (u *User) start(add func(*User)) error {
	success := false
	for !success {
		u.Send("Do you have an account? Y/n/q")
		command, err := u.readBlank()
		if err != nil {
			return err
		}

		switch strings.ToLower(command[0]) {
		case "n":
			err = u.createLogin()
			success = true
		case "q":
		case "qq":
		case "quit":
			if quit, err := u.quit(); quit || err != nil {
				return err
			}
		default:
			success, err = u.login()
		}

		if err != nil {
			return err
		}
	}

	add(u)
	u.Send("Welcome " + u.name)
	return nil
}

func (u *User) read() ([]string, error) {
	select {
	case <-u.Ctx.Done():
		return nil, errors.New("user closed")
	case text := <-u.RecieveChan:
		text = strings.Trim(text, " ")
		if len(text) == 0 {
			return nil, errors.New("no input detected")
		}
		return strings.Split(text, " "), nil
	}
}

func (u *User) readBlank() ([]string, error) {
	select {
	case <-u.Ctx.Done():
		return nil, errors.New("user closed")
	case text := <-u.RecieveChan:
		return strings.Split(strings.Trim(text, " "), " "), nil
	}
}

//Send sends a message to the user
func (u *User) Send(message string) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		defer close(out)
		select {
		case <-u.Ctx.Done():
		default:
			u.SendChan <- message
		}
	}()
	return out
}

func (u *User) kill(message string) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		defer close(out)
		select {
		case <-u.Ctx.Done():
		case <-u.Send(message):
			u.Cancel()
		}
	}()
	return out
}
