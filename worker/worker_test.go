package main

import (
	"fmt"
	"testing"

	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {

	tests := []struct {
		subscriber  *kafka.SubscriberMock
		expected    []byte
		expectedErr error
	}{
		{
			subscriber: &kafka.SubscriberMock{
				Message: []byte(`oi eu sou goku`),
				Err:     nil,
			},
			expected:    []byte(`oi eu sou goku`),
			expectedErr: nil,
		},
		{
			subscriber: &kafka.SubscriberMock{
				Message: nil,
				Err:     fmt.Errorf("some error on subscriber"),
			},
			expected:    nil,
			expectedErr: fmt.Errorf("some error on subscriber"),
		},
	}

	for _, test := range tests {
		m, err := readMessage(test.subscriber)

		assert.Equal(t, test.expected, m)
		if err != nil {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}

	}
}

func TestPublishMessage(t *testing.T) {
	tests := []struct {
		publisher   *kafka.PublishMock
		expectedErr error
	}{
		{
			publisher: &kafka.PublishMock{
				Err: nil,
			},
			expectedErr: nil,
		},
		{
			publisher: &kafka.PublishMock{
				Err: fmt.Errorf("Failed to publish message"),
			},
			expectedErr: fmt.Errorf("Failed to publish message"),
		},
	}
	nfe := &kafka.DBMessage{}

	for _, test := range tests {
		err := publishMessage(test.publisher, nfe)

		assert.Equal(t, test.expectedErr, err)

	}
}

var (
	correctMessageWithProcNFE = []byte(`
	{
		"access_key": "uhuuul",
		"xml": "PD94bWwgdmVyc2lvbj0iMS4wIj8+CjxuZmVQcm9jIHhtbG5zPSJodHRwOi8vd3d3LnBvcnRhbGZpc2NhbC5pbmYuYnIvbmZlIiB2ZXJzYW89IjMuMTAiPgoJPE5GZSB4bWxucz0iaHR0cDovL3d3dy5wb3J0YWxmaXNjYWwuaW5mLmJyL25mZSI+CiAgICAgICAgICAgICAgPGluZk5GZT4KCQk8dG90YWw+CgkJCTxJQ01TVG90PgoJCQkJPHZORj4xODAuNTA8L3ZORj4KCQkJPC9JQ01TVG90PgoJCTwvdG90YWw+CiAgICAgICAgICAgIDwvaW5mTkZlPgoJPC9ORmU+CjwvbmZlUHJvYz4="
	}
	`)

	correctMessageWithNFE = []byte(`
	{
		"access_key": "uhuuul",
		"xml": "PD94bWwgdmVyc2lvbj0iMS4wIj8+Cgk8TkZlIHhtbG5zPSJodHRwOi8vd3d3LnBvcnRhbGZpc2NhbC5pbmYuYnIvbmZlIj4KICAgICAgICAgICAgICA8aW5mTkZlPgoJCTx0b3RhbD4KCQkJPElDTVNUb3Q+CgkJCQk8dk5GPjE4MC41MDwvdk5GPgoJCQk8L0lDTVNUb3Q+CgkJPC90b3RhbD4KICAgICAgICAgICAgPC9pbmZORmU+Cgk8L05GZT4="
	}
	`)
	incorrectMessageWithProcNFE = []byte(`
	{
		"access_key": "uhuuul",
		"xml": "bmZlUHJvYyB2YyB0YSB0ZW50YW5kbyBkZWNvZGlmaWNhciBhIG1lbnNzYWdlPw=="
	}
	`)

	incorrectMessageWithNFE = []byte(`
	{
		"access_key": "uhuuul",
		"xml": "dmMgdGEgdGVudGFuZG8gZGVjb2RpZmljYXIgYSBtZW5zc2FnZT8="
	}
	`)
)

func TestProcessMessage(t *testing.T) {
	tests := []struct {
		expectedDbMessage *kafka.DBMessage
		message           []byte
		expectedErr       error
	}{
		{
			expectedDbMessage: &kafka.DBMessage{
				AccessKey: "uhuuul",
				Amount:    "180.50",
			},
			message:     correctMessageWithProcNFE,
			expectedErr: nil,
		},
		{
			expectedDbMessage: &kafka.DBMessage{
				AccessKey: "uhuuul",
				Amount:    "180.50",
			},
			message:     correctMessageWithNFE,
			expectedErr: nil,
		},
		{
			expectedDbMessage: nil,
			message:           incorrectMessageWithProcNFE,
			expectedErr:       fmt.Errorf("EOF"),
		},
		{
			expectedDbMessage: nil,
			message:           incorrectMessageWithNFE,
			expectedErr:       fmt.Errorf("EOF"),
		},
	}

	for _, test := range tests {
		dbMessage, err := processMessage(test.message)

		assert.Equal(t, test.expectedDbMessage, dbMessage)
		if err != nil {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}

	}
}
