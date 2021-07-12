package room

import "github.com/gorilla/websocket"
import "github.com/Pallinder/go-randomdata"

type Connection struct {
	Conn *websocket.Conn
	Room *Room
	Name string
}

func NewConnection(conn *websocket.Conn, room *Room) *Connection {
	connection := &Connection{
		Conn: conn,
		Room: room,
		Name: randomdata.SillyName(),
	}

	// send nickname to client
	conn.WriteMessage(
		websocket.TextMessage,
		[]byte("nickname " + connection.Name),
	)

	return connection
}
