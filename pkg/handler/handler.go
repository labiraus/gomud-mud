package handler

import (
	"fmt"
	"net/http"

	"github.com/labiraus/gomud-mud/pkg/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Start begins the api
func Start(g *game.Game) {
	http.HandleFunc("/game", userChan{
		userConnections: g.Users,
		ctx:             g.Ctx,
	}.gameHandler)
	fmt.Println("starting game server")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// gameHandler handles connections to the /game endpoint
func (users userChan) gameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	send := make(chan string, 10)
	recieve := make(chan string, 10)
	user := game.User{
		RecieveChan: recieve,
		SendChan:    send,
	}
	user.Init(users.ctx)

	// Send messages until user is closed
	go func() {
		for {
			select {
			case message := <-send:
				err := conn.WriteMessage(1, []byte(message))
				if err != nil {
					fmt.Println(err)
				}
			case <-user.Ctx.Done():
				select {
				case message := <-send:
					err := conn.WriteMessage(1, []byte(message))
					if err != nil {
						fmt.Println(err)
					}
				default:
					err := conn.WriteMessage(1, []byte("Disconnected."))
					if err != nil {
						fmt.Println(err)
					}
					conn.Close()
					return
				}
			}
		}
	}()

	// Read messages until closed
	go func() {
		var in []byte
		for {
			if _, in, err = conn.ReadMessage(); err != nil {
				fmt.Println(err.Error())
				user.Cancel()
				return
			}

			select {
			case <-user.Ctx.Done():
				fmt.Println("out!")
				return
			default:
				recieve <- string(in)
			}
		}
	}()

	users.userConnections <- &user
}
