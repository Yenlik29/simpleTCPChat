package main

import (
	"io"
	"time"
)

func NewWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{writer: writer,}
}

func (w *CommandWriter) WriteString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}, client *client) error {
	var name string

	if command.(MessageCommand).Name == "" {
		name = "unknown"
	} else {
		name = command.(MessageCommand).Name[:len(command.(MessageCommand).Name)]
	}

	curTime := time.Now()
	timing := curTime.Format("2006-01-02 15:04:05")
	msg := "\n[" + timing + "][" + name + "]:" + command.(MessageCommand).Message + "[" + timing + "][" + client.name + "]:"
	if err := w.WriteString(msg); err != nil {
		return err
	}
	return nil
}
