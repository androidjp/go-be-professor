package es

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestDocumentForES8 struct {
	TextID  int    `json:"text_id"`
	UID     string `json:"uid"`
	Content string `json:"content"`
}

func (t *TestDocumentForES8) ConvertBySearchHit(hit interface{}) error {
	if h, ok := hit.(types.Hit); ok {
		if err := sonic.Unmarshal(h.Source_, &t); err != nil {
			return err
		}
	}
	return nil
}

func (t *TestDocumentForES8) Clone() Document {
	return &TestDocumentForES8{
		TextID:  t.TextID,
		UID:     t.UID,
		Content: t.Content,
	}
}

func (t *TestDocumentForES8) GetID() (string, error) {
	return fmt.Sprintf("%d_%s", t.TextID, t.UID), nil
}

func TestClientES8_CreateIndex(t *testing.T) {

	Convey("Given a valid Elasticsearch index and mapping rule", t, func() {
		index := TestIndexMaExpert
		// 定义索引映射（包含IK分词器配置）
		mapping := `{
		"mappings": {
			"properties": {
				 	"content": {
						"type": "text",
						"analyzer": "ik_max_word",
						"search_analyzer": "ik_max_word"
					}
				}
			}
		}`

		Convey("When creating an index", func() {
			cli, _ := GetTestCfg().BuildES8Client()
			err := cli.CreateIndex(context.Background(), index, mapping)
			So(err, ShouldBeNil)
		})

	})
}

func TestClientES8_UpsertDocument(t *testing.T) {
	Convey("Given a valid Elasticsearch index and document", t, func() {
		cli, _ := GetTestCfg().BuildES8Client()

		Convey("When upserting a document", func() {
			doc := &TestDocumentForES8{
				TextID:  6,
				UID:     "user1",
				Content: "amy",
			}

			Convey("Then the document should be upserted successfully", func() {
				err := cli.UpsertDocument(context.Background(), TestIndexMaExpert, doc)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestClientES8_SearchDocumentsRaw(t *testing.T) {
	Convey("Given a valid Elasticsearch index and keyword", t, func() {
		cli, _ := GetTestCfg().BuildES8Client()

		Convey("When searching for documents, key=content, matchInfo=2223333", func() {
			index := TestIndexMaExpert
			raw := `{"query": {"match": {"content": "2223333"}}}`
			//raw := `{"match": {"content": "2223333"}}`
			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocumentForES8{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocumentForES8).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocumentForES8).Content, ShouldEqual, "2223333")
			})
		})
		Convey("When searching for documents, key=content, matchInfo=hello world", func() {
			index := TestIndexMaExpert
			raw := `{"query": {"match": {"content": "hello world"}}}`
			//raw := `{"match": {"content": "hello world"}}`
			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocumentForES8{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocumentForES8).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocumentForES8).Content, ShouldEqual, "hello world")
			})
		})
	})
}
