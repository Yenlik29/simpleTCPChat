package main

import (
	"os"
	"fmt"
	"sync"

	_ "github.com/marcusolsson/tui-go"
)

func NewServer() *TCPChat {
	return &TCPChat{
		quantity: 0,
		mutex: &sync.Mutex{},
	}
}

func main() {
	var Chat *TCPChat

	port := ":8989"
	Chat = NewServer()
	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else {
		fmt.Println("[USAGE]: ./simpleTCPChat $port")
		return
	}
	if err := Chat.Listen(port); err != nil {
		fmt.Println("[ERROR]: Failed to connect to TCP")
		return
	} else {
		Chat.Start()
	}
}
