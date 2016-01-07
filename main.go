package main

import (
	"os"
)

func main() {
	projectId := os.Getenv("GCLOUD_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	logglyUrl := os.Getenv("LOGGLY_URL")
	logInfo(LogData{
		"project_id":   projectId,
		"pubsub_topic": topicName,
		"loggly_url":   logglyUrl,
	})

	msgC := consumePubSubMsgs(projectId, topicName)
	postToLoggly(logglyUrl, msgC)
}
