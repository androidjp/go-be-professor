package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
)

var (
	ErrES7NotSupportKnnSearch = errors.New("es7 not support knnSearch")
)

type ClientES7 struct {
	es *elastic.Client
}

func (c *ClientES7) CreateIndex(ctx context.Context, index, indexMappingRule string) error {
	// 创建索引
	createIndex, err := c.es.CreateIndex(index). // 索引
							BodyString(indexMappingRule). // 索引映射
							Do(ctx)

	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("create index %s is not acknowledged", index)
	}
	return nil
}

func (c *ClientES7) DeleteIndex(ctx context.Context, index string) error {
	// 删除索引
	deleteIndex, err := c.es.DeleteIndex(index).
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
func (c *ClientES7) UpsertDocument(ctx context.Context, index string, doc Document) error {
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
	exists, err := c.es.Exists().
		Index(index).
		Id(docID).
		Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		// 如果文档存在，更新文档
		_, err = c.es.Update().Index(index).Id(docID).Doc(doc).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		// 如果文档不存在，插入文档
		idx, err := c.es.Index().Index(index).Id(docID).BodyJson(doc).Do(ctx)
		if err != nil {
			return err
		}
		fmt.Println(idx)
	}
	return nil
}

// DeleteDocument 删除文档
func (c *ClientES7) DeleteDocument(ctx context.Context, index string, doc Document) error {
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
	_, err = c.es.Delete().Index(index).Id(docID).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// SearchDocumentsMatchQuery 搜索文档（传递值的方式）
func (c *ClientES7) SearchDocumentsMatchQuery(ctx context.Context, index, key, matchInfo string, doc Document) ([]Document, error) {
	query := elastic.NewMatchQuery(key, matchInfo)
	searchResult, err := c.es.Search().
		Index(index).
		Query(query).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// 处理搜索结果
	var results []Document
	var item Document
	for _, hit := range searchResult.Hits.Hits {
		item = doc.Clone()
		if err = item.ConvertBySearchHit(hit); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

// SearchDocumentsRaw 搜索文档(原始ES SQL方式)
func (c *ClientES7) SearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error) {
	searchResult, err := c.es.Search().
		Index(index).
		Source(raw).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// 处理搜索结果
	var results []Document
	var item Document
	for _, hit := range searchResult.Hits.Hits {
		item = doc.Clone()
		if err = item.ConvertBySearchHit(hit); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (c *ClientES7) KnnSearchDocumentsRaw(ctx context.Context, index, raw string, doc Document) ([]Document, error) {
	return nil, ErrES7NotSupportKnnSearch
}

func (c *ClientES7) DeleteDocumentsRaw(ctx context.Context, index, raw string) error {
	_, err := c.es.DeleteByQuery(index).Q(raw).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
