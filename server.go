package main

import (
	"io"
	"net"
	"fmt"
	"time"
	"errors"
	"strings"
	"reflect"
)

func (Chat *TCPChat) Remove(client *client) {
	Chat.mutex.Lock()
	defer Chat.mutex.Unlock()

	name := client.name
	for i, check := range Chat.clients {
		if check == client {
			Chat.clients = append(Chat.clients[:i], Chat.clients[i+1:]...)
		}
	}

	fmt.Println("=====Closing connection=====")
	client.conn.Close()

	Chat.quantity--
	for _, client := range Chat.clients {
		curTime := time.Now()
		timing := curTime.Format("2006-01-02 15:04:05")
		msg := "\n[LEFT]:" + name + "\n[" + timing + "][" + client.name + "]:"
		client.writer.writer.Write([]byte(msg))
	}
}

func (Chat *TCPChat) BroadCast(command interface{}, client *client) error {
	for _, client := range Chat.clients {
		if client.name != command.(MessageCommand).Name {
			blank := strings.TrimSpace(command.(MessageCommand).Message) == ""
			if !blank{
				client.writer.Write(command, client)
			}
		} else if client.name == command.(MessageCommand).Name {
			client.Prefix()
		}
		fmt.Println("===", client.name, command.(MessageCommand).Name)
	}
	return nil
}

func (Chat *TCPChat) SetName(command string, client *client) error {
	client.name = command[:len(command)-1]
	for _, user := range Chat.clients {
		if user.name == "" {
			continue
		} else if user.name != client.name {
			curTime := time.Now()
			timing := curTime.Format("2006-01-02 15:04:05")
			msg := "\n" + client.name + " has joined our chat...\n[" + timing + "]" + "[" + user.name + "]:"
			user.writer.writer.Write([]byte(msg))
		}
	}
	client.Prefix()
	return nil
}


func (Chat *TCPChat) Accept(Conn net.Conn) (*client, error) {
	fmt.Println("=====NEW CONNECTION=====")
	fmt.Println("[ADDR]:", Conn.RemoteAddr().String())

	if Chat.quantity == 10 {
		Conn.Write([]byte("Chat is full of connections\n"))
		Conn.Close()
		return nil, errors.New("Maximum connections reached.")
	}

	Chat.mutex.Lock()
	defer Chat.mutex.Unlock()

	client := &client{conn: Conn, writer: NewWriter(Conn),}
	Chat.clients = append(Chat.clients, client)
	return client, nil
}

func (Chat *TCPChat) Start() {
	for {
		conn, err := Chat.listener.Accept()
		if err != nil {
			fmt.Println("[ERROR]: Failed to accept connection:", err)
		} else {
			client, err := Chat.Accept(conn)
			if err == nil {
				Chat.quantity++
				go Chat.Serve(client)
			} else {
				fmt.Println("[ERROR]:", err)
			}
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

func (client *client) WelcomeMessage() error {
	msg := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n"
	msg = msg + "[ENTER YOUR NAME]: "
	if err := client.writer.WriteString(msg); err != nil {
		return err
	}
	return nil
}

func (client *client) Prefix() error {
	curTime := time.Now()
	timing := curTime.Format("2006-01-02 15:04:05")
	msg := "[" + timing + "]" + "[" + client.name + "]:"
	if err := client.writer.WriteString(msg); err != nil {
		return err
	}
	return nil
}

func (Chat *TCPChat) Serve(client *client) {
	var count int

	count = 0
	cmd := NewReader(client.conn)
	defer Chat.Remove(client)

	client.WelcomeMessage()
	for {
		msg, err := cmd.Read(count)
		if err != nil && err != io.EOF {
			fmt.Println("[ERROR]:", err)
		}
		if msg != nil {
			if reflect.TypeOf(msg).String() == "main.SendCommand" && count != 0 {
				go Chat.BroadCast(MessageCommand{Name: client.name, Message: msg.(SendCommand).Message}, client)
			} else if reflect.TypeOf(msg).String() == "main.SendCommand" && count == 0 {
				go Chat.SetName(msg.(SendCommand).Message, client)
			}
		}
		if err == io.EOF {
			break
		}
		count++
	}
}
