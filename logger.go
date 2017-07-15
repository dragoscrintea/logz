package main

import (
	"encoding/json"
	"io"
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

func (l *logger) Info(event string, d data) {
	d["type"] = "info"
	l.writeJSON(event, d)
}

func (l *logger) Action(event string, d data) {
	d["type"] = "action"
	l.writeJSON(event, d)
}

func (l *logger) writeJSON(event string, d data) {
	d["event"] = event
	d["timestamp"] = time.Now().Format(time.RFC3339Nano)
	for k, v := range l.global {
		d[k] = v
	}
	l.writer.Encode(d)
}
