package main

import (
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
)

func consumePubSubMsgs(projectId string, topicName string) <-chan []byte {
	msgC := make(chan []byte, 100)
	subName := topicName + "-2loggly"
	ackDeadline := time.Second * 20

	client, err := google.DefaultClient(context.Background(), pubsub.ScopePubSub)
	if err != nil {
		logFatal(LogData{
			"event": "pubsub.error.auth",
			"error": err.Error(),
		})
	}
	ctx := cloud.NewContext(projectId, client)

	isSub, err := pubsub.SubExists(ctx, subName)
	if err != nil {
		logFatal(LogData{
			"event": "pubsub.error.subexists",
			"error": err.Error(),
		})
	}
	if !isSub {
		err = pubsub.CreateSub(ctx, subName, topicName, ackDeadline, "")
		if err != nil {
			logFatal(LogData{
				"event": "pubsub.error.createsub",
				"error": err.Error(),
			})
		}
	}

	go func() {
		for {
			msgs, err := pubsub.PullWait(ctx, subName, 20)
			if err != nil {
				logInfo(LogData{
					"event": "pubsub.error.pull",
					"error": err.Error(),
				})
				time.Sleep(time.Second)
				continue
			}

			for _, m := range msgs {
				msgC <- m.Data
				err = pubsub.Ack(ctx, subName, m.AckID)
				if err != nil {
					logInfo(LogData{
						"event": "pubsub.error.ack",
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return msgC
}
