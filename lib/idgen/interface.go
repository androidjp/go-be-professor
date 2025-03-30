package idgen

import "context"

type IdGenerator interface {
	NextId(ctx context.Context) (int64, error)
}

type SeqGenerator interface {
	NextSeq(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)
}
