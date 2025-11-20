package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)



func SetupJetstreamConfig()(jetstream.StreamConfig, error){
	streamName := os.Getenv("JS_STREAM_NAME")
	if streamName == "" {
		return jetstream.StreamConfig{}, errors.New("env var JS_STREAM_NAME is required")
	}
	streamSubject := os.Getenv("JS_STREAM_SUBJECTS") + os.Getenv("SOURCE")
	if streamSubject == "" {
		return jetstream.StreamConfig{}, errors.New("env var JS_STREAM_SUBJECT is required")
	}
	retentionHours, err := strconv.Atoi(os.Getenv("JS_RETENTION_HOURS"))
	if err != nil {
		return jetstream.StreamConfig{}, errors.New("env var JS_RETENTION_HOURS must be a valid integer")
	}

	jetstreamConfig := jetstream.StreamConfig {
		Name: streamName,
		Subjects: []string{streamSubject},
		MaxAge: time.Duration(retentionHours) * time.Hour,
	}
	
	return jetstreamConfig, nil
}

