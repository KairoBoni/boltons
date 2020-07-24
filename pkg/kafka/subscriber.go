package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

//Subscriber Structure that knows how to read a message in a topic
type Subscriber struct {
	reader *kafka.Reader
}

type SubscriberInterface interface {
	Read() ([]byte, error)
	Close()
}

//NewTopicSubscription create a new publisher on topic the specific topic
func NewTopicSubscription(kafkaBrokerUrls, kafkaClientId, kafkaTopic string) *Subscriber {
	brokers := strings.Split(kafkaBrokerUrls, ",")

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

//Read a new message that arrived in the topic
func (s *Subscriber) Read() ([]byte, error) {
	m, err := s.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	return m.Value, nil
}

//Close the publisher
func (s *Subscriber) Close() {
	s.reader.Close()
}
