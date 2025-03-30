package kafka

// import (
// 	"bytes"
// 	"encoding/json"

// 	"mylib/mq"
// 	"mylib/mq/cmn"

// 	// "ksogit.kingsoft.net/kgo/async"
// )

// type KgoKafkaConsumer struct {
// 	sub *async.Consumer
// }

// func (k *KgoKafkaConsumer) RegisterHandler(topic string, handler mq.MessageHandleFunc, errHandleFunc func(error)) {
// 	callback := func(headers mq.Headers, value []byte) error {
// 		msg := cmn.NewMessage(topic, value, cmn.WithHeaders(headers))
// 		err, ack := handler(msg)
// 		if err != nil {
// 			errHandleFunc(err)
// 			if !ack {
// 				return async.ErrRetry
// 			}
// 		}
// 		return nil
// 	}
// 	k.sub.RegisterHandlerV2(topic, callback)
// }

// func (k *KgoKafkaConsumer) Start() error {
// 	go k.sub.Consume()
// 	return nil
// }

// func (k *KgoKafkaConsumer) Stop() error {
// 	k.sub.Close()
// 	return nil
// }

// func (k *KgoKafkaConsumer) SetEncoder(encoder async.Encoder) {
// 	k.sub.SetEncoder(encoder)
// }

// type jsonEncoder struct{}

// func (e *jsonEncoder) Encode(v interface{}) ([]byte, error) {
// 	return json.Marshal(v)
// }

// func (e *jsonEncoder) Decode(data []byte, v interface{}) error {
// 	d := json.NewDecoder(bytes.NewReader(data))
// 	d.UseNumber()
// 	return d.Decode(v)
// }

// func NewKgoKafkaConsumer(cfg *KgoKafkaConsumerConfig) mq.Consumer {
// 	asyncCfg := &async.Config{
// 		ConsumerProject:      cfg.Project,
// 		ConsumerGroup:        cfg.GetGroup(),
// 		WorkersEachPartition: cfg.WorkersEachPartition,
// 		MaxRetries:           cfg.MaxRetries,
// 	}
// 	sub := async.NewKafkaConsumer(cfg.Domain, asyncCfg, cfg.Addr)
// 	sub.SetEncoder(&jsonEncoder{})
// 	return &KgoKafkaConsumer{
// 		sub: sub,
// 	}
// }
