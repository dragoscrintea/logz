package main

import (
	"os"
	"time"
)

func main() {
	projectId := os.Getenv("GCLOUD_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	logglyUrl := os.Getenv("LOGGLY_URL")

	msgC := consumePubSubMsgs(projectId, topicName)
	go postToLoggly(logglyUrl, msgC)

	time.Sleep(time.Second * 10)
}
