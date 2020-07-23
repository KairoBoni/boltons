package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	reader *kafka.Reader
}

func NewTopicSubscription(kafkaBrokerUrls, kafkaClientId, kafkaTopic string) *Subscriber {
	brokers := strings.Split(kafkaBrokerUrls, ",")

	// make a new reader that consumes from topic-A
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	return &Subscriber{
		reader: kafka.NewReader(config),
	}
}

func (s *Subscriber) Read() ([]byte, error) {
	m, err := s.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	return m.Value, nil
}

func (s *Subscriber) Close() {
	s.reader.Close()
}
