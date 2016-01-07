package main

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"
)

func postToLoggly(url string, logC <-chan []byte) {
	var wg sync.WaitGroup
	stopping := false

	client := &http.Client{}
	ticker := time.NewTicker(2 * time.Second)

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

				req, err := http.NewRequest("POST", url, &data)
				if err != nil {
					log.Println(err)
					return
				}
				req.Header.Set("Content-Type", "text/plain")

				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Println("Loggly POST failed with status code ", resp.StatusCode)
					return
				}

			}(buffer)
		}

		if stopping {
			ticker.Stop()
			wg.Wait()
			return
		}
	}
}
