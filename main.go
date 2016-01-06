package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
)

func main() {
	projectId := os.Getenv("GCLOUD_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	logglyUrl := os.Getenv("LOGGLY_URL")
	subName := topicName + "-2loggly"
	ackDeadline := time.Second * 20

	client, err := google.DefaultClient(context.Background(), pubsub.ScopePubSub)
	if err != nil {
		log.Fatal("Failed to create oAuth2 HTTP client", err)
	}
	ctx := cloud.NewContext(projectId, client)

	isSub, err := pubsub.SubExists(ctx, subName)
	if err != nil {
		log.Fatal("Failed to determine subscription state", err)
	}
	if !isSub {
		err = pubsub.CreateSub(ctx, subName, topicName, ackDeadline, "")
		if err != nil {
			log.Fatal("Failed to create subscription", err)
		}
	}

	logC := make(chan []byte, 100)
	go batchPoster(logglyUrl, logC)

	for {
		msgs, err := pubsub.PullWait(ctx, subName, 20)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		for _, m := range msgs {
			logC <- m.Data
			err = pubsub.Ack(ctx, subName, m.AckID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func batchPoster(url string, logC chan []byte) {
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
