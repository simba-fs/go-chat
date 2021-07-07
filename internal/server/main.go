package server

import (
	"fmt"
	"log"
	"net/http"
	"errors"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// handler for home page
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello world</h1>")
}

// handler for websocket echo 
func echo(w http.ResponseWriter, r *http.Request) {
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

	conn.SetCloseHandler(func(code int, text string) error{
		fmt.Printf("code = %d text = %s\n", code, text)
		return errors.New("this is a close error")
	})

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// Listen server on `addr(string)`
func Listen(addr string){
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)

	log.Printf("Listen on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
