package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
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
		log.Error().Msgf("Missing env variables")
	}

	time.Sleep(time.Second * 30)

	s := kafka.NewTopicSubscription(kafkaBrokerURLs, kafkaClientID, kafkaSubscribeTopic)
	p := kafka.NewPublisherOnTopic(kafkaBrokerURLs, kafkaClientID, kafkaPublisherTopic)

	defer s.Close()
	for {
		fmt.Println("waitng message")
		m, err := s.Read()
		fmt.Println("recived")
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
		}

		nfe, err := processNFE(m)
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
		}

		publishMessage, err := json.Marshal(nfe)
		if err != nil {
			log.Error().Msgf("error to marshal message: %s", err.Error())
		}

		p.Publish(context.Background(), nil, publishMessage)

	}
}

var r = regexp.MustCompile("(nfeProc)")

func processNFE(data []byte) (*kafka.DBMessage, error) {
	var (
		amount string
		m      = &Message{}
	)

	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}

	dataNFE, err := base64.StdEncoding.DecodeString(m.XML)
	if err != nil {
		return nil, err
	}

	if r.Match(dataNFE) {
		nfe := &NFEProc{}
		err = xml.Unmarshal(dataNFE, nfe)
		if err != nil {
			return nil, err
		}
		amount = nfe.Amount
	} else {
		nfe := &NFE{}
		err = xml.Unmarshal(dataNFE, nfe)
		if err != nil {
			return nil, err
		}
		amount = nfe.Amount
	}

	return &kafka.DBMessage{
		AccessKey: m.AccessKey,
		Amount:    amount,
	}, nil
}
