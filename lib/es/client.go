package es

import (
	"context"
)

type Client interface {
	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index, indexMappingRule string) error
	// DeleteIndex 删除索引
	DeleteIndex(ctx context.Context, index string) error
	// UpsertDocument 检查文档是否存在，如果存在则更新，否则插入
	UpsertDocument(ctx context.Context, index string, doc Document) error
	// DeleteDocument 删除文档
	DeleteDocument(ctx context.Context, index string, doc Document) error
	// SearchDocumentsRaw 搜索文档(原始json方式)
	// index: 索引名称
	// raw: 搜索条件
	// raw 例子：
	// {"query": {"match": {"content": "2223333"}}}
	// from: 起始位置，默认为 0
	// size: 返回数量，不传则没有限制
	SearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error)
	// KnnSearchDocumentsRaw 向量方式搜索文档(原始json方式)
	KnnSearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error)
	// DeleteDocumentsRaw 删除文档(原始json方式)
	DeleteDocumentsRaw(ctx context.Context, index, raw string) error
}
