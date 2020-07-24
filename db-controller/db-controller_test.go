package main

import (
	"fmt"
	"testing"

	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

var (
	correctFormat = []byte(`
	{
		"access_key": "haha1212",
		"total": "120 dols"
	}
	`)

	badFormat = []byte(`
	{
		"access_key: "haha1212",
		"total": "120 dols"
	}
	`)
)

func TestMessageSender(t *testing.T) {

	tests := []struct {
		subscriber  kafka.SubscriberMock
		accessKey   string
		total       string
		expectedErr error
	}{
		{
			subscriber: kafka.SubscriberMock{
				Message: correctFormat,
				Err:     nil,
			},
			accessKey:   "haha1212",
			total:       "120 dols",
			expectedErr: nil,
		},
		{
			subscriber: kafka.SubscriberMock{
				Message: badFormat,
				Err:     nil,
			},
			accessKey:   "",
			total:       "",
			expectedErr: fmt.Errorf("invalid character 'h' after object key"),
		},
		{
			subscriber: kafka.SubscriberMock{
				Message: correctFormat,
				Err:     fmt.Errorf("failed to read subscriped topic"),
			},
			accessKey:   "",
			total:       "",
			expectedErr: fmt.Errorf("failed to read subscriped topic"),
		},
	}

	for _, test := range tests {
		accessKey, total, err := reciveMessage(&test.subscriber)

		assert.Equal(t, test.accessKey, accessKey)
		assert.Equal(t, test.total, total)
		if err != nil {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}

	}
}
