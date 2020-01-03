package main 

import (
	"fmt"
	"sync"
)

func NewServer() *TCPChat {
	return &TCPChat{
		mutex: &sync.Mutex{},
	}
}

func main() {
	var Chat *TCPChat

	Chat = NewServer()
	if err := Chat.Listen(":2525"); err != nil {
		fmt.Println("[ERROR]: Failed to connect to TCP")
		return
	} else {
		Chat.Start()
	}
}
