package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	server   string
	nickname string
	room     string
}

var Prompt = ">"

func GetPrompt(client Client) string {
	return fmt.Sprintf("%s %s ", client.server, Prompt)
}

// ReadCommand reads command from stdin
func ReadCommand(client Client) ([]string, string, error) {
	// prompt
	fmt.Print(GetPrompt(client))

	// read
	reader := bufio.NewReader(os.Stdin)
	raw, err := reader.ReadString('\n')
	if len(raw) <= 1 {
		return []string{""}, "", nil
	}

	// remove '\n'
	raw = raw[:len(raw)-1]

	if err != nil {
		return nil, "", err
	}

	cmd := []string{}
	for _, val := range strings.Split(raw, " ") {
		if val != "" {
			cmd = append(cmd, val)
		}
	}

	return cmd, raw, nil
}

func Send(client Client, data string) {
	if client.conn == nil {
		return
	}
	fmt.Printf("send: %s\n", data)
	err := client.conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		fmt.Println("sending error")
	}
}

func main() {
	var client = Client{}

	for {
		cmd, raw, _ := ReadCommand(client)

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
				continue
			}
			client = Client{c, addr, "testUser", "testRoom"}

		case "/disconnect":
			err := client.conn.Close()
			if err != nil {
				fmt.Println("disconnect error")
			}
			client = Client{nil, "", "", ""}
		case "/room":
		case "/nickname":
		default:
			Send(client, raw)
		}
	}
}
