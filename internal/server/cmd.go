package server

import "fmt"
import "errors"
import "github.com/simba-fs/go-chat/internal/room"
import "github.com/simba-fs/go-chat/internal/cmdParser"
import "github.com/gorilla/websocket"

var cl cmdParser.CmdList

var (
	ErrNoConnection = errors.New("no connection")
	ErrNoRoom = errors.New("no default room")
)

// func wrapper(f func()) func(raw string, cmds []string, exec cmdParser.FuncExec)(string, error){
//
// }

func init(){
	c := []cmdParser.Cmd{
		cmdParser.New("msg", "echo message", func(raw string, cmds []string, exec cmdParser.FuncExec, arg ...interface{})(string, error){
			fmt.Printf("cmdd %v\n", cmds)

			// get conn
			// conn, ok := arg[0].(*websocket.Conn)
			// if !ok {
			//     return "", ErrNoConnection
			// }

			// get room
			// curtRoom means current room
			curtRoom, ok := arg[1].(*room.Room)
			if !ok {
				return "", ErrNoRoom
			}

			curtRoom.Broadcast("msg", raw)

			// err := conn.WriteMessage(websocket.TextMessage, []byte(raw))
			return "", nil
		}),
	}

	cl = cmdParser.CmdList{
		Cmds: c,
		Help: "help",
		Helper: cmdParser.Helper,
	}
}

func exec(cmd []byte, conn *websocket.Conn, room *room.Room){
	fmt.Println(cl.Exec(string(cmd), conn, room))
}
