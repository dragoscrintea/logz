package main

import (
	"os"
)

func main() {
	log := newLogger(os.Stdout, data{"service": "log.shipper"})

	projectID := os.Getenv("PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	uploadURL := os.Getenv("UPLOAD_URL")

	log.Info(data{
		"project_id":   projectID,
		"pubsub_topic": topicName,
		"upload_url":   uploadUrl,
	})

	msgC := consumePubSubMsgs(projectID, topicName)
	postToLoggly(uploadURL, msgC)
}
