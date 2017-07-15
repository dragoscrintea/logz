package main

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

func postToLoggly(url string, logC <-chan []byte) {
	var wg sync.WaitGroup
	stopping := false

	client := &http.Client{}
	ticker := time.NewTicker(5 * time.Second)

	for {
		var buffer bytes.Buffer

	buffering:
		for {
			select {
			case logLine, more := <-logC:
				if !more {
					stopping = true
					break buffering
				}
				if buffer.Len() > 0 {
					buffer.Write([]byte("\n"))
				}
				buffer.Write(logLine)

			case <-ticker.C:
				if buffer.Len() > 0 {
					break buffering
				}
			}
		}

		if buffer.Len() > 0 {
			wg.Add(1)
			go func(data bytes.Buffer) {
				defer wg.Done()
				start := time.Now()

				req, err := http.NewRequest("POST", url, &data)
				if err != nil {
					logInfo(logData{
						"event": "request.error.create",
						"error": err.Error(),
					})
					return
				}
				req.Header.Set("Content-Type", "text/plain")

				resp, err := client.Do(req)
				if err != nil {
					logInfo(logData{
						"event": "request.error.execute",
						"error": err.Error(),
					})
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					logInfo(logData{
						"event":       "request.error.response",
						"status_code": resp.StatusCode,
					})
					return
				}

				logInfo(logData{
					"event": "request.ok",
					"timer": time.Now().Sub(start).Nanoseconds(),
					"size":  buffer.Len(),
				})

			}(buffer)
		}

		if stopping {
			ticker.Stop()
			wg.Wait()
			return
		}
	}
}
