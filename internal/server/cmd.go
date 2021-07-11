package server

import "fmt"
import "errors"
import "github.com/simba-fs/go-chat/internal/cmdParser"
import "github.com/gorilla/websocket"

var cl cmdParser.CmdList

var (
	ErrNoConnection = errors.New("no connection")
)

// func wrapper(f func()) func(raw string, cmds []string, exec cmdParser.FuncExec)(string, error){
//
// }

func init(){
	c := []cmdParser.Cmd{
		cmdParser.New("msg", "echo message", func(raw string, cmds []string, arg interface{}, exec cmdParser.FuncExec)(string, error){
			fmt.Printf("cmdd %v\n", cmds)
			conn, ok := arg.(*websocket.Conn)
			if !ok {
				return "", ErrNoConnection
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(raw))
			return "", err
		}),
	}

	cl = cmdParser.CmdList{
		Cmds: c,
		Help: "help",
		Helper: cmdParser.Helper,
	}
}

func exec(cmd []byte, conn *websocket.Conn){
	fmt.Println(cl.Exec(string(cmd), conn))
}
