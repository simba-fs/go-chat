package room

import "github.com/gorilla/websocket"

type Room struct {
	Name string
	Conns []*websocket.Conn
}

// Broadcast send message to all connection in the room
func (room *Room) Broadcast(msgType string, data string) *Room {
	for _, i := range room.Conns {
		i.WriteMessage(websocket.TextMessage, []byte(data))
	}
	return room
}

// Add add connection to the room
func (room *Room) Add(conn *websocket.Conn) *Room {
	for _, i := range room.Conns {
		if i == conn {
			return room
		}
	}
	room.Conns = append(room.Conns, conn)
	return room
}

// Remove remove connection from the room
func (room *Room) Remove(conn *websocket.Conn) *Room {
	result := []*websocket.Conn{}
	for _, i := range room.Conns {
		if i != conn {
			result = append(result, i)
		}
	}
	room.Conns = result
	return room
}

// Clear clear all connection from the room
func (room *Room) Clear() *Room {
	room.Conns = []*websocket.Conn{}
	return room
}

func (room *Room) size() int {
	return len(room.Conns)
}

// New create a new empty room
func New(name string) *Room{
	room := &Room{name, []*websocket.Conn{}}
	return room
}
