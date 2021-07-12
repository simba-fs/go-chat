package room

import "github.com/gorilla/websocket"

type Room struct {
	Name string
	Conns []*Connection
}

// Broadcast send message to all Connection in the room
func (room *Room) Broadcast(msgType string, data string) *Room {
	for _, i := range room.Conns {
		i.Conn.WriteMessage(websocket.TextMessage, []byte(data))
	}
	return room
}

// Add add Connection to the room
func (room *Room) Add(conn *Connection) *Room {
	for _, i := range room.Conns {
		if i == conn {
			return room
		}
	}
	room.Conns = append(room.Conns, conn)
	return room
}

// Remove remove Connection from the room
func (room *Room) Remove(conn *Connection) *Room {
	result := []*Connection{}
	for _, i := range room.Conns {
		if i != conn {
			result = append(result, i)
		}
	}
	room.Conns = result
	return room
}

// Clear clear all Connection from the room
func (room *Room) Clear() *Room {
	room.Conns = []*Connection{}
	return room
}

func (room *Room) size() int {
	return len(room.Conns)
}

var rooms = []*Room{}

// Get returns a Room by the given name. If it doesn't exist, it will create a new one
func Get(name string) *Room{
	for _, i := range rooms {
		if i.Name == name {
			return i
		}
	}
	room := &Room{name, []*Connection{}}
	rooms = append(rooms, room)
	return room
}
