package main

import (
	"fmt"
	"flag"

	"github.com/simba-fs/go-chat/internal/server"
	"github.com/simba-fs/go-chat/internal/client"
)

func main(){
	port := flag.Int("port", 3000, "port to use")
	host := flag.String("host", "0.0.0.0", "ip to use")
	isServer := flag.Bool("server", false, "server mode")
	flag.Parse();

	addr := fmt.Sprintf("%s:%d", *host, *port)

	if(*isServer){
		server.Listen(addr)
	}else{
		client.Start()
	}
}
