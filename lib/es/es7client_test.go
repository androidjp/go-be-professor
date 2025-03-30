package es

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	TestIndexMaExpert               = "test_index_ma_expert"
	TestIndexMaExpertVectorIndexing = "test_index_ma_expert_vector"
)

func TestClient_CreateIndex(t *testing.T) {

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
			cli, _ := GetTestCfg().BuildES7Client()
			err := cli.CreateIndex(context.Background(), index, mapping)
			So(err, ShouldBeNil)
		})

	})
}

type TestDocument struct {
	TextID  int    `json:"text_id"`
	UID     string `json:"uid"`
	Content string `json:"content"`
}

func (t *TestDocument) ConvertBySearchHit(hit interface{}) error {
	if h, ok := hit.(*elastic.SearchHit); ok {
		if err := sonic.Unmarshal(h.Source, &t); err != nil {
			return err
		}
	}
	return nil
}

func (t *TestDocument) Clone() Document {
	return &TestDocument{
		TextID:  t.TextID,
		UID:     t.UID,
		Content: t.Content,
	}
}

func (t *TestDocument) GetID() (string, error) {
	return fmt.Sprintf("%d_%s", t.TextID, t.UID), nil
}

func TestClient_UpsertDocument(t *testing.T) {
	Convey("Given a valid Elasticsearch index and document", t, func() {
		cli, _ := GetTestCfg().BuildES7Client()

		Convey("When upserting a document", func() {
			doc := &TestDocument{
				TextID:  7,
				UID:     "user1",
				Content: "你好，小花，明天和后天去小明家吧",
			}

			Convey("Then the document should be upserted successfully", func() {
				err := cli.UpsertDocument(context.Background(), TestIndexMaExpert, doc)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestConfig_SearchDocumentsMatchQuery(t *testing.T) {
	Convey("Given a valid Elasticsearch index and keyword", t, func() {
		cli, _ := GetTestCfg().BuildES7Client()

		Convey("When searching for documents, key=content, matchInfo=2223333", func() {
			index := TestIndexMaExpert
			key := "content"
			matchInfo := "2223333"

			results, err := cli.SearchDocumentsMatchQuery(context.Background(), index, key, matchInfo, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "2223333")
			})
		})
		Convey("When searching for documents, key=content, matchInfo=hello world", func() {
			index := TestIndexMaExpert
			key := "content"
			matchInfo := "hello world"

			results, err := cli.SearchDocumentsMatchQuery(context.Background(), index, key, matchInfo, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "hello world")
			})
		})
	})
}

func TestConfig_SearchDocumentsRaw(t *testing.T) {
	Convey("Given a valid Elasticsearch index and keyword", t, func() {
		cli, _ := GetTestCfg().BuildES7Client()

		Convey("When searching for documents, key=content, matchInfo=2223333", func() {
			index := TestIndexMaExpert
			//raw := `{"match": {"content": "2223333"}}`
			raw := `{
			   "query": {
			       "match": {
			           "content": "2223333"
			       }
			   }
			}`
			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "2223333")
			})
		})
		Convey("When searching for documents, key=content, matchInfo=hello world", func() {
			index := TestIndexMaExpert
			//raw := `{"match": {"content": "hello world"}}`
			raw := `{
			   "query": {
			       "match": {
			           "content": "hello world"
			       }
			   }
			}`
			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "hello world")
			})
		})

		Convey("When searching for documents, key=content, matchInfo=world, 采用倒排索引", func() {
			index := TestIndexMaExpert
			raw := `{
		"query": {
			"match": {
				"content": {
					"query": "world",
					"operator": "and",
					"minimum_should_match": "5%"
				}
			}
		},
		"highlight": {
			"fields": {
				"content": {
					"number_of_fragments": 2,
					"fragment_size": 150
				}
			},
			"pre_tags": ["<em>"],
			"post_tags": ["</em>"]
		}
	}`

			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "hello world")
			})
		})
		Convey("When searching for documents, key=content, matchInfo=明天, 采用倒排索引", func() {
			index := TestIndexMaExpert
			raw := `{
		"query": {
			"match": {
				"content": {
					"query": "小花",
					"operator": "and",
					"minimum_should_match": "50%"
				}
			}
		},
		"highlight": {
			"fields": {
				"content": {
					"number_of_fragments": 2,
					"fragment_size": 150
				}
			},
			"pre_tags": ["<em>"],
			"post_tags": ["</em>"]
		}
	}`

			results, err := cli.SearchDocumentsRaw(context.Background(), index, raw, &TestDocument{})

			Convey("Then the results should match the expected documents", func() {
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].(*TestDocument).UID, ShouldEqual, "user1")
				So(results[0].(*TestDocument).Content, ShouldEqual, "你好，小花，明天和后天去小明家吧")
			})
		})
	})
}
