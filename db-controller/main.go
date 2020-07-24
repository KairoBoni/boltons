package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

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

	//Wait for a while until the Kafka and Postgres start
	time.Sleep(time.Second * 30)

	s := kafka.NewTopicSubscription(kafkaBrokerURLs, kafkaClientID, kafkaTopic)
	defer s.Close()

	store, err := database.NewStore(dbConfigFilepath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new store")
	}

	for {
		accessKey, total, err := reciveMessage(s)
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
		}

		if err := store.InsertNfeTotal(accessKey, total); err != nil {
			log.Error().Msgf("error while save database: %s", err.Error())
		}
	}
}

func reciveMessage(s kafka.SubscriberInterface) (string, string, error) {
	var NFE = &kafka.DBMessage{}

	m, err := s.Read()
	if err != nil {
		return "", "", err
	}

	if err := json.Unmarshal(m, NFE); err != nil {
		return "", "", err
	}

	return NFE.AccessKey, NFE.Total, nil

}
