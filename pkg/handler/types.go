package handler

import (
	"context"

	"pkg/game"
)

type userChan struct {
	userConnections chan<- *game.User
	ctx             context.Context
}

type UserRequest struct {
	UserName string
}

func (r UserRequest) Validate() error {
	return nil
}

type UserResponse struct {
	Greeting string
}

func (r UserResponse) Validate() error {
	return nil
}
