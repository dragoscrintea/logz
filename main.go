package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeout = time.Hour * 36

// Responds to environment variables:
//   PROJECT_ID (no default)
//   PUBSUB_TOPIC (no default)
//   UPLOAD_URL (no default)
func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		exitCode = 1
	}
}

func run() error {
	// parse env vars
	projectID := os.Getenv("PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	uploadURL := os.Getenv("UPLOAD_URL")

	// setup this app logging
	log := newLogger(os.Stdout, data{"service": "log.shipper"})
	log.Info("daemon.start", data{
		"project_id":   projectID,
		"pubsub_topic": topicName,
		"upload_url":   uploadURL,
	})

	// create daemon context
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		select {
		case <-sigC:
			log.Info("daemon.interrupted", data{})
			shutdown()
		case <-ctx.Done():
		}
	}()

	// open connection to pubsub
	t, err := newTopic(ctx, projectID, topicName)
	if err != nil {
		log.Action("topic.failed", data{"error": err.Error()})
		shutdown()
		return err
	}
	err = t.Subscribe(ctx, func(logLine []byte, err error) {
		if err != nil {
			// TODO: beware log explosion
		} else {
			// TODO: add logline
		}
	})
	t.Stop()
	if err != nil {
		log.Action("subscribe.failed", data{"error": err.Error()})
		shutdown()
		return err
	}

	// msgC := consumePubSubMsgs(projectID, topicName)
	// postToLoggly(uploadURL, msgC)

	return nil
}
