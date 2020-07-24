package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

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
		messageChan         = make(chan kafka.WorkerMessage)
		errChan             = make(chan error)
		doneChan            = make(chan bool)
		goRoutines          = 3
	)
	defer close(messageChan)
	defer close(errChan)
	defer close(doneChan)

	if kafkaBrokerURLs == "" || kafkaTopic == "" || kafkaClientID == "" || arquiveiCredentials == "" {
		log.Fatal().Msgf("Missing some environment variable (CREDENTIALS_FILEPATH|KAFKA_BROKERS|KAFKA_CLIENT_ID|KAFKA_TOPIC)")
	}

	//Wait for a while until the Kafka and Postgres start
	time.Sleep(time.Second * 30)

	cli, err := NewArquiveiClient(arquiveiCredentials)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new arquivei client")
	}

	p := kafka.NewPublisherOnTopic(kafkaBrokerURLs, kafkaClientID, kafkaTopic)
	defer p.Close()

	for i := 0; i < goRoutines; i++ {
		go publishMessage(messageChan, p, errChan)
	}

	go func() {
		for {
			NFEs, err := cli.RequestNFEs()
			if err != nil {
				errChan <- fmt.Errorf("Failed to request NFE: %v", err)
			}
			if len(NFEs) < 1 {
				doneChan <- true
			}
			for _, workerMessage := range NFEs {
				messageChan <- workerMessage
			}
		}
	}()

	for {
		select {
		case err := <-errChan:
			log.Fatal().Err(err).Msgf("Failed to send message to worker")

		case <-doneChan:
			log.Info().Msg("All NFEs were successfully collected")
			return

		}
	}
}

func publishMessage(messageChan chan kafka.WorkerMessage, p kafka.PublisherInterface, errChan chan error) {
	for {
		m, ok := (<-messageChan)
		if !ok {
			return
		}

		message, err := json.Marshal(m)
		if err != nil {
			errChan <- fmt.Errorf("Failed to marshal message: %v", err)
		}

		err = p.Publish(context.Background(), nil, message)
		if err != nil {
			errChan <- fmt.Errorf("Failed to publish mensssage: %v", err)
		}
	}
}
