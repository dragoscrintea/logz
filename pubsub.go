package main

import (
	"log"
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

	go func() {
		for {
			msgs, err := pubsub.PullWait(ctx, subName, 20)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}

			for _, m := range msgs {
				msgC <- m.Data
				err = pubsub.Ack(ctx, subName, m.AckID)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return msgC
}
