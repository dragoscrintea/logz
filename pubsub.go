package main

import (
	"encoding/json"
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
		var lastAckErr error
		var lastDecodeErr error

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

			lastAckErr = nil
			lastDecodeErr = nil

			for _, m := range msgs {
				jsonBytes, err := msgToLogJson(m.Data)
				if err != nil {
					lastDecodeErr = err
				} else {
					msgC <- jsonBytes
				}

				err = pubsub.Ack(ctx, subName, m.AckID)
				if err != nil {
					lastAckErr = err
				}
			}

			// protect against flood of logs from this application
			// when acks/decodes fail but adding rollup & sleep delay
			if lastAckErr != nil {
				logInfo(LogData{
					"event": "pubsub.error.ack",
					"error": lastAckErr.Error(),
				})
				time.Sleep(time.Second)
			}
			if lastDecodeErr != nil {
				logInfo(LogData{
					"event": "pubsub.error.decode",
					"error": lastDecodeErr.Error(),
				})
				time.Sleep(time.Second)
			}
		}
	}()

	return msgC
}

type cloudLogJson struct {
	Meta struct {
		Timestamp string `json:"timestamp"`
		Label     struct {
			PodName string `json:"container.googleapis.com/pod_name"`
		} `json:"labels"`
	} `json:"metadata"`
	Text    string          `json:"textPayload"`
	RawJson json.RawMessage `json:"structPayload"`
}

func msgToLogJson(msgData []byte) ([]byte, error) {
	log := cloudLogJson{}

	err := json.Unmarshal(msgData, &log)
	if err != nil {
		return []byte{}, err
	}

	var jsonBytes []byte
	if len(log.Text) > 0 {
		jsonBytes, err = json.Marshal(map[string]string{
			"timestamp": log.Meta.Timestamp,
			"pod":       log.Meta.Label.PodName,
			"text":      log.Text,
		})
	} else {
		jsonBytes, err = log.RawJson.MarshalJSON()
	}

	return jsonBytes, err
}
