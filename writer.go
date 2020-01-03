package main

import (
	"io"
)

func NewWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{writer: writer,}
}

func (w *CommandWriter) WriteString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	var name string

	if command.(MessageCommand).Name == "" {
		name = "unknown"
	} else {
		name = command.(MessageCommand).Name[:len(command.(MessageCommand).Name)-1]
	}

	msg := "[" + name + "]:" + command.(MessageCommand).Message
	if err := w.WriteString(msg); err != nil {
		return err
	}
	return nil
}
