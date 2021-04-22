package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/labiraus/gomud-common/db"
	"github.com/labiraus/gomud-mud/pkg/game"
	"github.com/labiraus/gomud-mud/pkg/handler"
)

// This example demonstrates a trivial echo server.
func main() {
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	ctx, ctxDone := context.WithCancel(context.Background())
	defer ctxDone()
	db.Setup(ctx, dbName, dbUser, dbPassword)
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
