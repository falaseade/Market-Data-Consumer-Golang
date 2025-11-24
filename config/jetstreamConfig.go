package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type JetStreamConfig struct {
	Name string
	Subject string
	Source string
	RetentionHours int
}


func LoadJetStreamConfigFromEnv() (JetStreamConfig, error) {
	streamName := os.Getenv("JS_STREAM_NAME")
	if streamName == "" {
		return JetStreamConfig{}, fmt.Errorf("JS_STREAM_NAME stream name must be set")
	}
	streamSubject := os.Getenv("JS_STREAM_SUBJECTS")
	if streamSubject == "" {
		return JetStreamConfig{}, fmt.Errorf("JS_STREAM_SUBJECTS stream subject must be set")
	}
	source := os.Getenv("SOURCE")
	if source == "" {
		return JetStreamConfig{}, fmt.Errorf("SOURCE must be set")
	}
	retentionHours, err := strconv.Atoi(os.Getenv("JS_RETENTION_HOURS"))
	if err != nil || retentionHours <= 0 {
		return JetStreamConfig{}, fmt.Errorf("JS_RETENTION_HOURS stream retention time not set or less than zero")
	}
	return JetStreamConfig{
		Name: streamName,
		Subject: streamSubject,
		Source: source,
		RetentionHours: retentionHours,
	}, nil
}


func (j JetStreamConfig) ToStreamConfig()(jetstream.StreamConfig, error){
	if j.Name == "" {
		return jetstream.StreamConfig{}, fmt.Errorf("JS_STREAM_NAME stream name must be set")
	}

	if j.Subject == "" {
		return jetstream.StreamConfig{}, fmt.Errorf("JS_STREAM_SUBJECTS stream subject must be set")
	}

	if j.Source == "" {
		return jetstream.StreamConfig{}, fmt.Errorf("SOURCE must be set")
	}

	if j.RetentionHours <= 0 {
		return jetstream.StreamConfig{}, fmt.Errorf("JS_RETENTION_HOURS stream retention time must be set")
	}

	return jetstream.StreamConfig{
		Name: j.Name,
		Subjects: []string{fmt.Sprintf("%s.%s", j.Subject, j.Source)},
		MaxAge: time.Duration(j.RetentionHours) * time.Hour,
	}, nil
}



