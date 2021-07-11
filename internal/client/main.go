package client

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gorilla/websocket"

	"github.com/simba-fs/go-chat/internal/cmdParser"
)

type Client struct {
	conn     *websocket.Conn
	server   string
	nickname string
	room     string
}

var client = &Client{nil, "", "", ""}

func livePrompt() (string, bool) {
	if client.conn == nil {
		return "> ", true
	}
	return fmt.Sprintf("%s#%s @ %s > ", client.nickname, client.room, client.server), true
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "/connect", Description: "connect to chat server"},
		{Text: "/disconnect", Description: "disconnect"},
		{Text: "/room", Description: "enter/list a chat room"},
		{Text: "/nickname", Description: "change/show nickname"},
		{Text: "/help", Description: "show this message"},
		{Text: "/exit", Description: "exit this program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Send(msgType string, data string) {
	if client.conn == nil {
		return
	}
	msg := fmt.Sprintf("%s %s", msgType, data)
	fmt.Printf("send: %s\n", msg)
	err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Println("sending error")
	}
}

func receive() {
	for {
		if client.conn == nil {
			continue
		}
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			break
		}
		cmd := strings.Split(string(msg), " ")
		if len(cmd) < 1 {
			continue
		}
		switch cmd[0] {
		case "msg":
			fmt.Printf("receive: %s\n", msg)
		case "room":
			client.room = cmd[1]
		case "nickname":
			client.nickname = cmd[1]
		}
	}
}

func Start() {
	c := []cmdParser.Cmd{
		cmdParser.New("/connect", "connect to chat server", func(raw string, cmds []string, exec cmdParser.FuncExec) (string, error) {
			addr := "ws://127.0.0.1:3000/echo"
			if len(cmds) > 1 {
				addr = cmds[1]
			}
			c, _, err := websocket.DefaultDialer.Dial(addr, nil)

			if err != nil {
				fmt.Println("connecting error")
			} else {
				client.conn = c
				client.server = addr
				client.nickname = "testUser"
				client.room = "testRoom"
			}
			return "", nil
		}),
		cmdParser.New("/disconnect", "disconnect", func(raw string, cmds []string, exec cmdParser.FuncExec) (string, error) {
			client.conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection close by client"),
				// NOTE: I don't know what the usage of deadline is, I don't need it here
				time.Now().AddDate(0, 0, 1),
			)
			client.conn = nil
			client.server = ""
			client.nickname = ""
			client.room = ""
			return "", nil
		}),
		cmdParser.New("/room [id]", "enter/list a chat room on server", func(raw string, cmds []string, exec cmdParser.FuncExec) (string, error) {
			if len(cmds) > 1 {
				Send("room", cmds[1])
				// on success
				client.room = cmds[1]
			} else {
				fmt.Printf("Your are at room %s\n", client.room)
			}
			return "", nil
		}),
		cmdParser.New("/nickname [name]", "change/show nickname", func(raw string, cmds []string, exec cmdParser.FuncExec) (string, error) {
			if len(cmds) > 1 {
				Send("nickname", cmds[1])
				// on success
				client.nickname = cmds[1]
			} else {
				fmt.Printf("Your nickname is %s\n", client.nickname)
			}
			return "", nil
		}),
		cmdParser.New("/exit", "exit program", func(raw string, cmds []string, exec cmdParser.FuncExec) (string, error) {
			exec("/disconnect")
			fmt.Println("Good bye ~")
			os.Exit(0)
			return "", nil
		}),
	}

	cl := cmdParser.CmdList{
		Cmds: c,
		Help: "/help",
		Helper: func(cmds *cmdParser.CmdList) string {
			fmt.Print(cmdParser.Helper(cmds))
			return ""
		},
	}

	go receive()
	prompt.New(
		func(raw string) {
			_, err := cl.Exec(raw)
			if err != nil {
				if err == cmdParser.ErrCommandNotFound {
					Send("msg", raw)
				}
			}
		},
		completer,
		prompt.OptionLivePrefix(livePrompt),
	).Run()
}
