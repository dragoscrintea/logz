package main

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

const (
	uploadEveryXSeconds  = 5
	httpTimeoutInSeconds = 30
)

func upload(log *logger, url string, logC <-chan []byte) {
	var wg sync.WaitGroup
	stopping := false

	client := &http.Client{
		Timeout: httpTimeoutInSeconds * time.Second,
	}
	ticker := time.NewTicker(uploadEveryXSeconds * time.Second)

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
			go func(buf bytes.Buffer) {
				defer wg.Done()
				start := time.Now()

				req, err := http.NewRequest("POST", url, &buf)
				if err != nil {
					log.Action("request.error.create", data{
						"error": err.Error(),
					})
					return
				}
				req.Header.Set("Content-Type", "text/plain")

				resp, err := client.Do(req)
				if err != nil {
					log.Action("request.error.execute", data{
						"error": err.Error(),
						"timer": time.Now().Sub(start).Seconds(),
						"size":  buffer.Len(),
					})
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Action("request.error.response", data{
						"status_code": resp.StatusCode,
						"timer":       time.Now().Sub(start).Seconds(),
						"size":        buffer.Len(),
					})
					return
				}

				log.Info("request.ok", data{
					"timer": time.Now().Sub(start).Seconds(),
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
