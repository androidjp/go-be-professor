package cmn

import "mylib/mq"

type Message struct {
	topic   string
	value   []byte
	headers mq.Headers
}

type MessageOption func(*Message)

func WithHeaders(headers mq.Headers) MessageOption {
	return func(msg *Message) {
		msg.headers = headers
	}
}

func NewMessage(topic string, value []byte, opts ...MessageOption) *Message {
	entry := &Message{
		topic:   topic,
		value:   value,
		headers: mq.NewHeaders(),
	}
	for _, opt := range opts {
		opt(entry)
	}
	return entry
}

func (m *Message) Headers() mq.Headers {
	return m.headers
}

func (m *Message) Value() []byte {
	return m.value
}

func (m *Message) Topic() string {
	return m.topic
}
