package es

import (
	"context"
	"fmt"
)

var ErrEsNotImplemented = fmt.Errorf("es not implemented")

type EmptyClient struct {
}

func (e EmptyClient) DeleteDocumentsRaw(ctx context.Context, index, raw string) error {
	return ErrEsNotImplemented
}

func (e EmptyClient) CreateIndex(ctx context.Context, index, indexMappingRule string) error {
	return ErrEsNotImplemented
}

func (e EmptyClient) DeleteIndex(ctx context.Context, index string) error {
	return ErrEsNotImplemented
}

func (e EmptyClient) UpsertDocument(ctx context.Context, index string, doc Document) error {
	return ErrEsNotImplemented
}

func (e EmptyClient) DeleteDocument(ctx context.Context, index string, doc Document) error {
	return ErrEsNotImplemented
}

func (e EmptyClient) SearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error) {
	return nil, ErrEsNotImplemented
}
func (e EmptyClient) KnnSearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error) {
	return nil, ErrEsNotImplemented
}
