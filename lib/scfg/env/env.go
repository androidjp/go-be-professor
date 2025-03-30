package env

import (
	"encoding/json"
	"mylib/scfg/kcfg"
	"os"
	"strings"
)

type Env struct {
	prefixes []string
}

func NewSource(prefix ...string) kcfg.Source {
	return &Env{
		prefixes: prefix,
	}
}

func (e *Env) Load() ([]*kcfg.KeyValue, error) {
	kv, err := e.load(os.Environ())
	if err != nil {
		return nil, err
	}
	return []*kcfg.KeyValue{kv}, nil
}

func (e *Env) Watch() (kcfg.Watcher, error) {
	return nil, nil
}

func (e *Env) Close() error {
	return nil
}

func (e *Env) SourceName() string {
	return "env"
}

func (e *Env) load(envs []string) (*kcfg.KeyValue, error) {
	kvMap := make(map[string]interface{})
	for _, env := range envs {
		var k, v string
		subs := strings.SplitN(env, "=", 2) //nolint:gomnd
		k = subs[0]
		if len(subs) > 1 {
			v = subs[1]
		}
		if len(e.prefixes) > 0 {
			p, ok := matchPrefix(e.prefixes, k)
			if !ok || len(p) == len(k) {
				continue
			}
		}
		if len(k) != 0 {
			kvMap[k] = v
		}
	}
	value, err := json.Marshal(kvMap)
	if err != nil {
		return nil, err
	}
	return &kcfg.KeyValue{
		Key:    "env_source",
		Value:  value,
		Format: "json",
	}, nil
}

func matchPrefix(prefixes []string, s string) (string, bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}
