package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

//Publisher Structure that knows how to publish a message in a topic
type Publisher struct {
	writer *kafka.Writer
}

type PublisherInterface interface {
	Publish(parent context.Context, key, value []byte) error
	Close()
}

//NewPublisherOnTopic create a new publisher on topic the specific topic
func NewPublisherOnTopic(kafkaBrokerUrls, clientId, topic string) *Publisher {
	brokers := strings.Split(kafkaBrokerUrls, ",")
	dialer := &kafka.Dialer{
		Timeout:  30 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return &Publisher{
		writer: kafka.NewWriter(config),
	}
}

//Publish send a new message
func (p *Publisher) Publish(parent context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(parent, message)
}

//Close the publisher
func (p *Publisher) Close() {
	p.writer.Close()
}
