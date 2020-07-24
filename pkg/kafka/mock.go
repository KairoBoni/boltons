package kafka

import "context"

//PublishMock use for test to another packages
//this struct implements the PublisherInterface
type PublishMock struct {
	Err error
}

//SubscriberMock use for test to another packages
//this struct implements the SubscriberInterface
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
