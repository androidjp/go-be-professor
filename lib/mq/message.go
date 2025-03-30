package mq

type Message interface {
	Headers() Headers
	Value() []byte
	Topic() string
}

type Headers map[string]interface{}

func NewHeaders() Headers {
	return make(map[string]interface{})
}
