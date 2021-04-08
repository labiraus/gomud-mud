package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"pkg/game"

	"github.com/labiraus/gomud-common/db"
	"github.com/labiraus/gomud-mud/pkg/handler"

	_ "github.com/denisenkom/go-mssqldb"
)

// This example demonstrates a trivial echo server.
func main() {
	fmt.Println("hi")
	ctx, ctxDone := context.WithCancel(context.Background())
	defer ctxDone()
	db.Setup(ctx)
	g, err := game.New(ctx)
	if err != nil {
		fmt.Println(err)
	}
	go handler.Start(g)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println("Got signal: " + s.String() + " now closing")
}
