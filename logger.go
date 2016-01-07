package main

import (
	"encoding/json"
	"os"
)

type LogData map[string]interface{}

func logInfo(data LogData) {
	os.Stdout.Write(map2Json(data))
	os.Stdout.Write([]byte("\n"))
}

func logFatal(data LogData) {
	os.Stdout.Write(map2Json(data))
	os.Stdout.Write([]byte("\n"))
	os.Exit(1)
}

func map2Json(d LogData) []byte {
	json, _ := json.Marshal(d)
	return json
}
