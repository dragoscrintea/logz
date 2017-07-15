package main

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type data map[string]interface{}

type logger struct {
	writer *json.Encoder
	global data
}

func newLogger(writer io.Writer, global data) *logger {
	return &logger{
		writer: json.NewEncoder(writer),
		global: global,
	}
}

func (l *logger) Info(d data) {
	d["type"] = "info"
	l.writeJSON(d)
}

func (l *logger) Fatal(d data) {
	d["type"] = "action"
	l.writeJSON(d)
	os.Exit(1)
}

func (l *logger) writeJSON(d data) {
	d["timestamp"] = time.Now().Format(time.RFC3339Nano)
	for k, v := range l.global {
		d[k] = v
	}
	l.writer.Encode(d)
}
