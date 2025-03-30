package kafka

// import (
// 	"context"
// 	"mylib/mq"

// 	"ksogit.kingsoft.net/kgo/async"
// )

// type KgoKafkaPublisher struct {
// 	pub *async.Producer
// }

// func NewKgoKafkaPublisher(cfg *KgoKafkaPublisherConfig) (*KgoKafkaPublisher, error) {
// 	pub, err := async.NewKafkaProducer(cfg.Domain, cfg.Addr, async.SetProject(cfg.Project))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &KgoKafkaPublisher{
// 		pub: pub,
// 	}, nil
// }

// func (p *KgoKafkaPublisher) Publish(ctx context.Context, msg mq.Message) error {
// 	_, err := p.pub.ProduceV2(msg.Topic(), msg.Headers(), msg.Value())
// 	return err
// }

// func (p *KgoKafkaPublisher) PublishWithPartitionKey(ctx context.Context, partitionKey uint64, msg mq.Message) error {
// 	_, err := p.pub.ProduceV3(msg.Topic(), partitionKey, msg.Headers(), msg.Value())
// 	return err
// }

// func (p *KgoKafkaPublisher) Release() {
// 	p.pub.Close()
// }
