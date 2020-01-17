package main

import (
	"io"
	"net"
	"fmt"
	"reflect"
)

func (Chat *TCPChat) Remove(client *client) {
	Chat.mutex.Lock()
	defer Chat.mutex.Unlock()

	for i, check := range Chat.clients {
		if check == client {
			Chat.clients = append(Chat.clients[:i], Chat.clients[i+1:]...)
		}
	}

	fmt.Println("=====Closing connection=====")
	client.conn.Close()
}

func (Chat *TCPChat) BroadCast(command interface{}) error {
	for _, client := range Chat.clients {
		if client.name != command.(MessageCommand).Name {
			client.writer.Write(command)
		}
		fmt.Println("===", client.name, command.(MessageCommand).Name)
		// client.writer.Write(command)
	}
	return nil
}

func (Chat *TCPChat) SetName(command string, client *client) error {
	client.name = command[:len(command)-1]
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

func (client *client) WelcomeMessage() error {
	msg := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--' hjm\n"
	msg = msg + "[ENTER YOUR NAME]: "
	if err := client.writer.WriteString(msg); err != nil {
		return err
	}
	return nil
}

func (client *client) Prefix() error {
	msg := "[MESSAGE]: "
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
		if count != 0 {
			client.Prefix()
		}
		fmt.Printf("[COUNT]:[%d][%s]\n", count, client.name)
		msg, err := cmd.Read(count)
		if err != nil && err != io.EOF {
			fmt.Println("[ERROR]:", err)
		}
		if msg != nil {
			if reflect.TypeOf(msg).String() == "main.SendCommand" && count != 0 {
				go Chat.BroadCast(MessageCommand{Name: client.name, Message: msg.(SendCommand).Message})
			} else if reflect.TypeOf(msg).String() == "main.SendCommand" && count == 0 {
				go Chat.SetName(msg.(SendCommand).Message, client)
			} else if reflect.TypeOf(msg).String() == "main.NameCommand" {
				client.name = msg.(NameCommand).Name
			}
		}
		if err == io.EOF {
			break
		}
		count++
	}
}
