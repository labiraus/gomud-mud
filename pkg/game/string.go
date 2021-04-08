package game

import (
	"context"

	"github.com/labiraus/gomud-common/core"
)

var globalActions = map[string]func(context.Context, core.Entity, []string) error{
	"qq":   quit,
	"quit": quit,
}
