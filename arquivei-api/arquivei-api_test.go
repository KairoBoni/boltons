package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/KairoBoni/boltons/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

var (
	normalResponse1 = []byte(`
	{
		"data": [
		{
			"access_key": "eae",
			"xml": "xml encodado"
		}
		],
		"page": {
			"next": "nextpage",
			"previous": "previouspage"
		}
	}
	`)
	normalResponse2 = []byte(`
	{
		"data": [
		{
			"access_key": "Pois bem",
			"xml": "cheguei"
		},
		{
			"access_key": "quero ficar",
			"xml": "bem avontade"
		},
		{
			"access_key": "na verdade",
			"xml": "sou assim"
		}
		],
		"page": {
			"next": "nextpage",
			"previous": "previouspage"
		}
	}
	`)

	thisResponseDoesNotMatter = []byte(`
	{
		blablalba
	}
	`)

	PoorlyFormattedResponse = []byte(`e se a respsota vier assim?`)
)

type MockArquiveiClient struct {
	resposeBody    []byte
	responseStatus int
	err            error
}

func (v *MockArquiveiClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader(v.resposeBody)),
		StatusCode: v.responseStatus,
	}, v.err
}

func TestRequestNFEs(t *testing.T) {
	tests := []struct {
		page        Page
		mock        MockArquiveiClient
		expected    []kafka.WorkerMessage
		expectedErr error
	}{
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    normalResponse1,
				responseStatus: http.StatusOK,
				err:            nil,
			},
			expected: []kafka.WorkerMessage{
				{
					AccessKey: "eae",
					XML:       "xml encodado",
				},
			},
			expectedErr: nil,
		},
		{
			page: Page{
				Next:     "e se vir outra pagina?",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    normalResponse2,
				responseStatus: http.StatusOK,
				err:            nil,
			},
			expected: []kafka.WorkerMessage{
				{
					AccessKey: "Pois bem",
					XML:       "cheguei",
				},
				{
					AccessKey: "quero ficar",
					XML:       "bem avontade",
				},
				{
					AccessKey: "na verdade",
					XML:       "sou assim",
				},
			},
			expectedErr: nil,
		},
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    thisResponseDoesNotMatter,
				responseStatus: http.StatusInternalServerError,
				err:            nil,
			},
			expected:    nil,
			expectedErr: fmt.Errorf("maximum retry attempts reached"),
		},
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    thisResponseDoesNotMatter,
				responseStatus: http.StatusUnauthorized,
				err:            nil,
			},
			expected:    nil,
			expectedErr: fmt.Errorf("user unauthorized"),
		},
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    thisResponseDoesNotMatter,
				responseStatus: http.StatusForbidden,
				err:            nil,
			},
			expected:    nil,
			expectedErr: fmt.Errorf("user forbidden"),
		},
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    thisResponseDoesNotMatter,
				responseStatus: http.StatusNotFound,
				err:            nil,
			},
			expected:    nil,
			expectedErr: fmt.Errorf("endpoint not found"),
		},
		{
			page: Page{
				Next:     "",
				Previous: "",
			},
			mock: MockArquiveiClient{
				resposeBody:    PoorlyFormattedResponse,
				responseStatus: http.StatusOK,
				err:            nil,
			},
			expected:    nil,
			expectedErr: fmt.Errorf("invalid character 'e' looking for beginning of value"),
		},
	}

	for _, test := range tests {
		cli := ArquiveiClient{
			client: &test.mock,
			page:   test.page,
			cred: &credentials{
				APIId:  "",
				APIKey: "",
			},
		}
		got, err := cli.RequestNFEs()
		assert.Equal(t, test.expected, got)
		if err != nil {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}

func TestMessageSender(t *testing.T) {
	tests := []struct {
		message     kafka.WorkerMessage
		publisher   kafka.PublishMock
		expectedErr error
	}{
		{
			message: kafka.WorkerMessage{
				AccessKey: "chave de acesso",
				XML:       "dlflasdsf==",
			},
			publisher: kafka.PublishMock{
				Err: nil,
			},
			expectedErr: nil,
		},
		{
			message: kafka.WorkerMessage{
				AccessKey: "chave de acesso",
				XML:       "dlflasdsf==",
			},
			publisher: kafka.PublishMock{
				Err: fmt.Errorf("Publish refused"),
			},
			expectedErr: fmt.Errorf("Failed to publish mensssage: Publish refused"),
		},
	}

	for _, test := range tests {

		messageChan := make(chan kafka.WorkerMessage)
		errChan := make(chan error)

		go publishMessage(messageChan, &test.publisher, errChan)

		messageChan <- test.message
		if test.expectedErr == nil {
			close(errChan)
		}
		close(messageChan)

		assert.Equal(t, test.expectedErr, <-errChan)
	}
}
