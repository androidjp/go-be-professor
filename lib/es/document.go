package es

type Document interface {
	GetID() (string, error)                   // 获取文档ID
	ConvertBySearchHit(hit interface{}) error // 从搜索结果中转换
	Clone() Document                          // 深拷贝的自定义实现
}
