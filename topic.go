package main

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

// Topic encapsulates the message queue topic
type topic struct {
	ProjectID string
	pubSub    *pubsub.Topic
}

// Stop cleanly closes the topic
func (t *topic) Stop() {
	if t.pubSub != nil {
		t.pubSub.Stop()
	}
}

// New creates a message queue topic
func newTopic(ctx context.Context, projectID string, topicName string) (*topic, error) {
	if len(projectID) == 0 {
		return &topic{}, nil
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	t := client.Topic(topicName)
	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		t, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, err
		}
	}

	return &topic{
		ProjectID: projectID,
		pubSub:    t,
	}, nil
}

// Subscribe reads JSON-encoded logs from the topic and decodes them
func (t *topic) Subscribe(ctx context.Context, fn func(s []byte, err error)) error {
	const ackDeadline = time.Second * 20

	if t.pubSub == nil {
		return errors.New("Topic has no project ID")
	}

	subName := t.pubSub.ID() + "." + time.Now().Format("v2006-01-02-15-04-05.999999")

	client, err := pubsub.NewClient(ctx, t.ProjectID)
	if err != nil {
		return err
	}

	sub, err := client.CreateSubscription(ctx, subName, t.pubSub, ackDeadline, nil)
	if err != nil {
		return err
	}
	defer sub.Delete(context.Background())

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fn(msgToLogJSON(m.Data))
		m.Ack()
	})
	if err != nil {
		return err
	}

	return nil
}

type cloudLogJSON struct {
	Meta struct {
		Timestamp string `json:"timestamp"`
		Label     struct {
			PodName string `json:"container.googleapis.com/pod_name"`
		} `json:"labels"`
	} `json:"metadata"`
	Text    string          `json:"textPayload"`
	RawJSON json.RawMessage `json:"structPayload"`
}

func msgToLogJSON(msgData []byte) ([]byte, error) {
	log := cloudLogJSON{}

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
		jsonBytes, err = log.RawJSON.MarshalJSON()
	}

	return jsonBytes, err
}
