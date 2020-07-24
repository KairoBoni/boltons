package kafka

import "context"

type PublishMock struct {
	Err error
}

type SubscriberMock struct {
	Message []byte
	Err     error
}

func (m *PublishMock) Publish(parent context.Context, key, value []byte) error {
	return m.Err
}

func (m *PublishMock) Close() {}

func (m *SubscriberMock) Read() ([]byte, error) {
	return m.Message, m.Err
}

func (m *SubscriberMock) Close() {}
