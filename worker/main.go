package main

import (
	"os"
	"time"

	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	var (
		kafkaBrokerURLs     = os.Getenv("KAFKA_BROKERS")
		kafkaClientID       = os.Getenv("KAFKA_CLIENT_ID")
		kafkaPublisherTopic = os.Getenv("KAFKA_PUBLISHER_TOPIC")
		kafkaSubscribeTopic = os.Getenv("KAFKA_SUBSCRIBE_TOPIC")
	)

	if kafkaBrokerURLs == "" || kafkaSubscribeTopic == "" || kafkaClientID == "" || kafkaPublisherTopic == "" {
		log.Error().Msgf("Missing some environment variable (KAFKA_BROKERS|KAFKA_CLIENT_ID|KAFKA_PUBLISHER_TOPIC|KAFKA_SUBSCRIBE_TOPIC)")
	}

	time.Sleep(time.Second * 30)

	s := kafka.NewTopicSubscription(kafkaBrokerURLs, kafkaClientID, kafkaSubscribeTopic)
	defer s.Close()
	p := kafka.NewPublisherOnTopic(kafkaBrokerURLs, kafkaClientID, kafkaPublisherTopic)
	defer p.Close()

	for {
		m, err := readMessage(s)
		if err != nil {
			log.Error().Msgf("Failed to receive message: %s", err.Error())
		}

		nfe, err := processMessage(m)
		if err != nil {
			log.Error().Msgf("Failed to process message: %s", err.Error())
		}

		if err = publishMessage(p, nfe); err != nil {
			log.Error().Msgf("Failed to publish message: %s", err.Error())
		}

	}
}
