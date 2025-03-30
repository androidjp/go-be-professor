package file

import (
	"fmt"
	"io"
	"mylib/scfg/kcfg"
	"os"
	"strings"
)

// LocalFile is a file source.
type LocalFile struct {
	path string
}

func NewSource(path string) kcfg.Source {
	return &LocalFile{path: path}
}

func (f *LocalFile) Load() ([]*kcfg.KeyValue, error) {
	fs, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	// 暂不支持读取目录
	if fs.IsDir() {
		return nil, fmt.Errorf("file source can not load file path is a dir")
	}

	kv, err := f.loadFile(f.path)
	return []*kcfg.KeyValue{kv}, err
}

func (f *LocalFile) loadFile(path string) (*kcfg.KeyValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &kcfg.KeyValue{
		Key:    info.Name(),
		Format: format(info.Name()),
		Value:  data,
	}, nil
}

func format(name string) string {
	if p := strings.Split(name, "."); len(p) > 1 {
		return p[len(p)-1]
	}
	return ""
}

// Watch is not supported.
func (f *LocalFile) Watch() (kcfg.Watcher, error) {
	return nil, nil
}

func (f *LocalFile) SourceName() string {
	return "local_file"
}

func (f *LocalFile) Close() error {
	return nil
}
