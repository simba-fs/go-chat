package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/simba-fs/go-chat/internal/room"
)

var upgrader = websocket.Upgrader{}
var rooms = []*room.Room{}
var defaultRoom = room.Get("default")

// handler for home page
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello world</h1>")
}

// handler for websocket echo
func wsServer(w http.ResponseWriter, r *http.Request) {
	// cros
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	connection := room.NewConnection(conn, defaultRoom)
	defaultRoom.Add(connection)

	conn.SetCloseHandler(func(code int, text string) error {
		return errors.New("this is a close error")
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		exec(message, connection, defaultRoom)
	}
}

// Listen server on `addr(string)`
func Listen(addr string) {
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", wsServer)

	log.Printf("Listen on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
