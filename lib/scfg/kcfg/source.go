package kcfg

// KeyValue : 代表这一个文件里所有内容
type KeyValue struct {
	Key    string // 文件名
	Value  []byte // 文件内容
	Format string // 文件格式：json、yaml、toml
}

// Source : 代表一个数据源，一个数据源会有多份配置文件
type Source interface {
	Load() ([]*KeyValue, error)
	Watch() (Watcher, error)
	Close() error
	SourceName() string
}

type Watcher interface {
	WaitRead() ([]*KeyValue, error)
	Stop() error
}
