package kafka

import "context"

type PublisherInterface interface {
	Publish(parent context.Context, key, value []byte) error
	Close()
}

type SubscriberInterface interface {
	Read() ([]byte, error)
	Close()
}

type WorkerMessage struct {
	AccessKey string `json:"access_key"`
	XML       string `json:"xml"`
}

type DBMessage struct {
	AccessKey string `json:"access_key"`
	Total     string `json:"total"`
}
