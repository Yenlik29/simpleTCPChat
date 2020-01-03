package main

import (
	"io"
	"net"
	"sync"
	"bufio"
)

type CommandReader struct {
	reader 		*bufio.Reader
}

type CommandWriter struct {
	writer 		io.Writer
}

type MessageCommand struct {
	Name 		string
	Message 	string
}

type SendCommand struct {
	Message 	string
}

type NameCommand struct {
	Name 		string
}

type client struct {
	conn 		net.Conn
	name 		string
	writer 		*CommandWriter
}

type TCPChat struct {
	listener 	net.Listener
	clients 	[]*client
	mutex 		*sync.Mutex
}
