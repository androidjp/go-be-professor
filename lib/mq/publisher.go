package mq

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, msg Message) error
	PublishWithPartitionKey(ctx context.Context, partitionKey uint64, msg Message) error
	Release()
}
