package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()
	var (
		dbConfigFilepath = os.Getenv("CONFIG_DB_FILEPATH")
		kafkaBrokerURLs  = os.Getenv("KAFKA_BROKERS")
		kafkaClientID    = os.Getenv("KAFKA_CLIENT_ID")
		kafkaTopic       = os.Getenv("KAFKA_TOPIC")
	)

	s := kafka.NewTopicSubscription(kafkaBrokerURLs, kafkaClientID, kafkaTopic)
	defer s.Close()

	store, err := database.NewStore(dbConfigFilepath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new store")
	}

	for {
		m, err := s.Read()

		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
		}

		totalNFE, err := decodeMessage(m)
		if err != nil {
			log.Error().Msgf("error while decoding message: %s", err.Error())
		}

		if err := store.InsertNfeTotal(totalNFE.AccessKey, totalNFE.Total); err != nil {
			log.Error().Msgf("error while save database: %s", err.Error())
		}
	}

}

func decodeMessage(message []byte) (*TotalNFE, error) {
	var totalNFE = &TotalNFE{}

	if err := json.Unmarshal(message, totalNFE); err != nil {
		return nil, err
	}

	return totalNFE, nil
}
