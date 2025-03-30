package es

import (
	"fmt"
)

const (
	// 定义索引映射（包含IK分词器配置）
	IK_MAX_WORD_MAPPING_RULE = `{
		"mappings": {
			"properties": {
				"%s": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				}
			}
		}
	}`
)

// GetIKMaxWordMappingRule 针对指定的某些字段，得到IK分词器最大分词规则
func GetIKMaxWordMappingRule(field string) string {
	if field == "" {
		return ""
	}
	return fmt.Sprintf(IK_MAX_WORD_MAPPING_RULE, field)
}
