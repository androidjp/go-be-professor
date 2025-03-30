package es

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type ClientES8 struct {
	es *elasticsearch.TypedClient
}

func (c *ClientES8) CreateIndex(ctx context.Context, index, indexMappingRule string) error {
	// 创建索引

	buf := bytes.NewBufferString(indexMappingRule)

	createIndex, err := c.es.Indices.Create(index).
		Raw(buf). // 索引映射
		Do(ctx)

	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("create index %s is not acknowledged", index)
	}
	return nil
}

func (c *ClientES8) DeleteIndex(ctx context.Context, index string) error {
	// 删除索引
	deleteIndex, err := c.es.Indices.Delete(index).
		Do(ctx)

	if err != nil {
		return err
	}
	if !deleteIndex.Acknowledged {
		return fmt.Errorf("delete index %s is not acknowledged", index)
	}
	return nil
}

// UpsertDocument 检查文档是否存在，如果存在则更新，否则插入
func (c *ClientES8) UpsertDocument(ctx context.Context, index string, doc Document) error {
	var (
		docID string
		err   error
	)

	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	// 得到文档ID
	docID, err = doc.GetID()
	if err != nil {
		return err
	}

	// 执行是否存在的查询
	exists, err := c.es.Exists(index, docID).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		// 如果文档存在，更新文档
		_, err = c.es.Update(index, docID).Doc(doc).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		// 如果文档不存在，插入文档
		idx, err := c.es.Index(index).Id(docID).Document(doc).Do(ctx)
		if err != nil {
			return err
		}
		fmt.Println(idx)
	}
	return nil
}

// DeleteDocument 删除文档
func (c *ClientES8) DeleteDocument(ctx context.Context, index string, doc Document) error {
	var (
		docID string
		err   error
	)

	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	// 得到文档ID
	docID, err = doc.GetID()
	if err != nil {
		return err
	}
	//删除文档
	_, err = c.es.Delete(index, docID).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// SearchDocumentsRaw 搜索文档(原始ES SQL方式)
func (c *ClientES8) SearchDocumentsRaw(ctx context.Context, index, query string, doc Document) ([]Document, error) {
	// 3. Search for the indexed documents
	// Init the request body.
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(query)

	// Perform the search request.
	res, err := c.es.Search().Index(index).Raw(buf).Do(ctx)

	if err != nil {
		return nil, err
	}

	// 处理搜索结果
	var results []Document
	var item Document
	for _, hit := range res.Hits.Hits {
		item = doc.Clone()
		if err = item.ConvertBySearchHit(hit); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

// KnnSearchDocumentsRaw 向量搜索文档(原始ES SQL方式)
func (c *ClientES8) KnnSearchDocumentsRaw(ctx context.Context, index, query string, doc Document) ([]Document, error) {
	// 3. Search for the indexed documents
	// Init the request body.
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(query)

	// Perform the search request.
	res, err := c.es.KnnSearch(index).Raw(buf).Do(ctx)

	if err != nil {
		return nil, err
	}

	// 处理搜索结果
	var results []Document
	var item Document
	for _, hit := range res.Hits.Hits {
		item = doc.Clone()
		if err = item.ConvertBySearchHit(hit); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (c *ClientES8) DeleteDocumentsRaw(ctx context.Context, index, raw string) error {
	// Init the request body.
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(raw)

	_, err := c.es.DeleteByQuery(index).Raw(buf).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
