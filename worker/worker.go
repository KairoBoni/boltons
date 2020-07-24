package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"regexp"

	"github.com/KairoBoni/boltons/pkg/kafka"
)

func readMessage(s kafka.SubscriberInterface) ([]byte, error) {
	m, err := s.Read()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func publishMessage(p kafka.PublisherInterface, nfe *kafka.DBMessage) error {
	message, err := json.Marshal(nfe)
	if err != nil {
		return err
	}
	if err = p.Publish(context.Background(), nil, message); err != nil {
		return err
	}
	return nil
}

var r = regexp.MustCompile("(nfeProc)")

//processMessage get the message of worker topic and transform in the message
//of db topic
func processMessage(message []byte) (*kafka.DBMessage, error) {
	var (
		amount string
		m      = &kafka.WorkerMessage{}
	)

	err := json.Unmarshal(message, m)
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
