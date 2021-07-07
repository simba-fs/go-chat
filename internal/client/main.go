package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	server   string
	nickname string
	room     string
}

var client = &Client{nil, "", "", ""}

func livePrompt() (string, bool) {
	return fmt.Sprintf("%s > ", client.server), true
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

// ReadCommand reads command from stdin
func parseCmd(raw string) []string {
	if len(raw) <= 0 {
		return []string{}
	}

	cmd := []string{}
	for _, val := range strings.Split(raw, " ") {
		if val != "" {
			cmd = append(cmd, val)
		}
	}

	return cmd
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

func executor(raw string) {
	cmd := parseCmd(raw)

	switch cmd[0] {
	case "/help", "/cmd", "":
		fmt.Println("/connect <URL>        connect to chat sercer")
		fmt.Println("/disconnect           disconnect")
		fmt.Println("/room [room id]       enter/list a chat room on server")
		fmt.Println("/nickname [new name]  change/show nickname")
		fmt.Println("/help, /cmd           show this message")
	case "/connect":
		addr := "ws://127.0.0.1:3000/echo"
		if len(cmd) > 1 {
			addr = cmd[1]
		}
		c, _, err := websocket.DefaultDialer.Dial(addr, nil)

		if err != nil {
			fmt.Println("connecting error")
		} else {
			client = &Client{c, addr, "testUser", "testRoom"}
		}

	case "/disconnect":
		client.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection close by client"),
			// NOTE: I don't know what the usage of deadline is, I don't need it here
			time.Now().AddDate(0, 0, 1),
		)
		client = &Client{nil, "", "", ""}
	case "/room":
		if len(cmd) > 1 {
			Send("room", cmd[1])
			// on success
			client.room = cmd[1]
		} else {
			fmt.Printf("Your are at room %s", client.room)
		}
	case "/nickname":
		if len(cmd) > 1 {
			Send("nickname", cmd[1])
			// on success
			client.nickname = cmd[1]
		} else {
			fmt.Printf("Your nickname is %s", client.nickname)
		}
	case "/exit":
		if client.conn != nil {
			client.conn.Close()
		}
		fmt.Println("Good bye ~")
		os.Exit(0)
	default:
		Send("msg", raw)
	}
}

func main() {
	prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(livePrompt),
	).Run()
}
