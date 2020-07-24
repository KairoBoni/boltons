package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisherOnTopic(kafkaBrokerUrls, clientId, topic string) *Publisher {
	brokers := strings.Split(kafkaBrokerUrls, ",")
	fmt.Println(brokers)
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

func (p *Publisher) Publish(parent context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(parent, message)
}

func (p *Publisher) Close() {
	p.writer.Close()
}
