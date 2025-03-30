package es

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetIKMaxWordMappingRule(t *testing.T) {
	Convey("Given a valid Elasticsearch index and mapping rule", t, func() {
		Convey("When creating an index", func() {
			mapping := GetIKMaxWordMappingRule("content1")
			So(mapping, ShouldNotBeEmpty)
			So(mapping, ShouldEqualJSON, `{
		"mappings": {
			"properties": {
				"content1": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				}
			}
		}
	}`)
		})
	})
}
