package main

import (
	"io"
	"net"
	"fmt"
	"reflect"
)

func (Chat *TCPChat) BroadCast(command interface{}) error {
	for _, client := range Chat.clients {
		client.writer.Write(command)
	}
	return nil
}

func (Chat *TCPChat) Accept(Conn net.Conn) *client {
	fmt.Println("=====NEW CONNECTION=====")
	fmt.Println("[ADDR]:", Conn.RemoteAddr().String())

	Chat.mutex.Lock()
	defer Chat.mutex.Unlock()

	client := &client{conn: Conn, writer: NewWriter(Conn),}
	Chat.clients = append(Chat.clients, client)
	return client
}

func (Chat *TCPChat) Start() {
	for {
		conn, err := Chat.listener.Accept()
		if err != nil {
			fmt.Println("[ERROR]: Failed to accept connection:", err)
		} else {
			client := Chat.Accept(conn)
			go Chat.Serve(client)
		}
	}
}

func (Chat *TCPChat) Listen(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("[ERROR]: Failed to listen to tcp")
		return err
	} else {
		Chat.listener = listen
		fmt.Println("[INFO]: Listening on the port", address)
		return nil
	}
}

func (Chat *TCPChat) Serve(client *client) {
	cmd := NewReader(client.conn)
	// defer Chat.Remove(client)
	for {
		cmdName, err := cmd.Read()
		if err != nil && err != io.EOF {
			fmt.Println("[ERROR]:", err)
		}
		if cmdName != nil {
			if reflect.TypeOf(cmdName).String() == "main.SendCommand" {
				go Chat.BroadCast(&MessageCommand{Name: client.name, Message: cmdName.(SendCommand).Message})
			} else if reflect.TypeOf(cmdName).String() == "main.NameCommand" {
				client.name = cmdName.(NameCommand).Name
			}
		}
		if err == io.EOF {
			break
		}
	}
}
