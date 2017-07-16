package main

import (
	"encoding/json"
	"errors"
	"strings"
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
func (t *topic) Subscribe(ctx context.Context, subName string, fn func(s []byte, err error)) error {
	const ackDeadline = time.Second * 20

	if t.pubSub == nil {
		return errors.New("Topic has no project ID")
	}

	client, err := pubsub.NewClient(ctx, t.ProjectID)
	if err != nil {
		return err
	}

	sub := client.Subscription(subName)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, subName, t.pubSub, ackDeadline, nil)
		if err != nil {
			return err
		}
	}

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fn(msgToLogJSON(m.Data))
		m.Ack()
	})
	if err != nil {
		return err
	}

	return nil
}

// Using Stackdriver v2 format
// Spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
// View actuals at: https://console.cloud.google.com/logs/viewer
type cloudLogJSON struct {
	Timestamp string `json:"timestamp"`
	Label     struct {
		PodName string `json:"container.googleapis.com/pod_name"`
	} `json:"labels"`
	Text    string          `json:"textPayload"`
	RawJSON json.RawMessage `json:"jsonPayload"`
}

func msgToLogJSON(msgData []byte) ([]byte, error) {
	log := cloudLogJSON{}

	err := json.Unmarshal(msgData, &log)
	if err != nil {
		return []byte{}, err
	}

	var jsonBytes []byte
	if len(log.Text) > 0 {
		// in stackdriver v2, the json detection has changed and raw JSON
		// is coming through as text, so temporarily:
		if strings.HasPrefix(log.Text, "{") {
			jsonBytes = []byte(strings.TrimSpace(log.Text))
		} else {
			jsonBytes, err = json.Marshal(map[string]string{
				"timestamp": log.Timestamp,
				"pod":       log.Label.PodName,
				"text":      log.Text,
			})
		}
	} else {
		jsonBytes, err = log.RawJSON.MarshalJSON()
	}

	return jsonBytes, err
}
