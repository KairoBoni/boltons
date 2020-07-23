package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()
	var (
		arquiveiCredentials = os.Getenv("CREDENTIALS_FILEPATH")
		kafkaBrokerURLs     = os.Getenv("KAFKA_BROKERS")
		kafkaClientID       = os.Getenv("KAFKA_CLIENT_ID")
		kafkaTopic          = os.Getenv("KAFKA_TOPIC")
	)

	if kafkaBrokerURLs == "" || kafkaTopic == "" || kafkaClientID == "" || arquiveiCredentials == "" {
		log.Fatal().Msgf("Missing env variables")
	}

	cli, err := NewArquiveiClient(arquiveiCredentials)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new arquivei client")
	}

	p := kafka.NewPublisherOnTopic(kafkaBrokerURLs, kafkaClientID, kafkaTopic)
	defer p.Close()

	for {
		NFCs, err := cli.RequestNFCs()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to request NFC")
		}
		if len(NFCs) < 1 {
			break
		}
		for _, nfc := range NFCs {
			message, err := json.Marshal(&nfc)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to umarshal")
			}
			err = p.Publish(context.Background(), nil, message)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to publisher")
			}
		}
	}

	log.Print("All NFCs are processed")

}
