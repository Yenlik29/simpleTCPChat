package main

import (
	"io"
	"fmt"
	"bufio"
	_ "errors"
)

func NewReader(reader io.Reader) *CommandReader {
	return &CommandReader{reader: bufio.NewReader(reader),}
}

func ReadSend(r *CommandReader) (interface{}, error) {
	msg, err := r.reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return SendCommand{}, err
	}

	fmt.Printf("[%s]\n", msg[:len(msg)-1])
	return SendCommand{Message: msg}, nil
}

func ReadName(r *CommandReader) (interface{}, error) {
	name, err := r.reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return NameCommand{}, err
	}

	fmt.Printf("[%s]\n", name[:len(name)-1])
	return NameCommand{Name: name}, nil
}

func (r *CommandReader) Read(count int) (interface{}, error) {
	cmdName, err := r.reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return nil, err
	}

	fmt.Printf("[CMD]:[%s]:", cmdName[:len(cmdName)-1])

	if count == 0 {
		return SendCommand{Message: cmdName}, nil
	} else {
		return SendCommand{Message: cmdName}, nil
	}
	return nil, nil
}
