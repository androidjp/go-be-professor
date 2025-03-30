package mq

type MessageHandleFunc func(message Message) (err error, ack bool)

type Consumer interface {
	RegisterHandler(topic string, handler MessageHandleFunc, errHandleFunc func(error))
	// SetEncoder(encoder async.Encoder)
	Start() error
	Stop() error
}
