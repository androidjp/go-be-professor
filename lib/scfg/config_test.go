package scfg_test

import (
	"fmt"
	"mylib/scfg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	path := "./../mock/cfg/env_and_test.yaml"

	envValue := "aaa"
	os.Setenv("KLION22_TEST_KEY", envValue)
	scfg.ParseConfig(path, nil)

	assert.Equal(t, true, scfg.GetBool("test.bool"))
	assert.Equal(t, "mylib", scfg.GetString("test.string"))
	assert.Equal(t, int32(-123), scfg.GetInt32("test.int"))
	assert.Equal(t, int64(-123), scfg.GetInt64("test.int"))
	assert.Equal(t, -123, scfg.GetInt("test.int"))
	assert.Equal(t, uint32(123), scfg.GetUint32("test.uint"))
	assert.Equal(t, uint64(123), scfg.GetUint64("test.uint"))
	assert.Equal(t, uint(123), scfg.GetUint("test.uint"))

	assert.Equal(t, 123.23, scfg.GetFloat64("test.float64"))

	assert.Equal(t, []string{"a", "b"}, scfg.GetStringSlice("test.stringSlice"))

	assert.Equal(t, map[string]string{"a": "b"}, scfg.GetStringMapString("test.mapStr"))

	expectStrMap := make(map[string]interface{})
	expectStrMap["a"] = 123
	assert.Equal(t, expectStrMap, scfg.GetStringMap("test.strMap"))
	assert.Equal(t, map[string][]string{"a": []string{"a", "b"}}, scfg.GetStringMapStringSlice("test.mapStrSlice"))

	missKey := "missKey"
	err := fmt.Errorf("key [%s] not found from [local_file] config source", missKey)
	assert.Equal(t, err, scfg.GetValue(missKey).Error())

	// test env
	assert.Equal(t, envValue, scfg.GetString("KLION22_TEST_KEY"))
}

func TestConfigApollo(t *testing.T) {
	path := "./../mock/cfg/app-apollo.yaml"
	scfg.ParseConfig(path, nil)

	// 测试用例： int、string、float、slice、json struct
	// intKey    => 123
	// strKey    => "string"
	// floatKey  => 234.5
	// sliceKey  => ["a", "b"]
	// jsonKey   => { "key": { "k1": "v1", "k2": "v2"}}
	type KV struct {
		K1 string `json:"k1"`
		K2 string `json:"k2"`
	}
	type jK struct {
		Key KV `json:"key"`
	}

	expectMap := make(map[string]interface{})
	expectMap["intKey"] = 123
	expectMap["strKey"] = "string"
	expectMap["floatKey"] = 234.5
	expectMap["sliceKey"] = []string{"a", "b"}
	expectMap["jsonKey"] = jK{
		Key: KV{
			K1: "v1",
			K2: "v2",
		},
	}

	for k, v := range expectMap {
		switch v.(type) {
		case int, int32, int64:
			require.Equal(t, v, scfg.GetInt(k))
		case string:
			require.Equal(t, v, scfg.GetString(k))
		case float64, float32:
			require.Equal(t, v, scfg.GetFloat64(k))
		case []string:
			require.Equal(t, v, scfg.GetStringSlice(k))
		default:
			actualVal := jK{}
			err := scfg.UnmarshalKey(k, &actualVal)
			require.Nil(t, err)
			require.Equal(t, v, actualVal)
		}

	}
}

// 具体看看这个说明文档： https://365.kdocs.cn/l/cqbqLYrKk9Xu
// 需要绑定 Host： 121.36.82.223  kacsvr-internal-beta.wps.cn 再进行调用
func TestConfigKAC(t *testing.T) {
	path := "./../mock/cfg/app-kac.yaml"
	scfg.ParseConfig(path, nil)

	// 测试用例： int、string、float、slice、json struct
	// intKey    => 123
	// strKey    => "string"
	// floatKey  => 234.5
	// sliceKey  => ["a", "b"]
	// jsonKey   => { "key": { "k1": "v1", "k2": "v2"}}
	type KV struct {
		K1 string `json:"k1"`
		K2 string `json:"k2"`
	}
	type jK struct {
		Key KV `json:"key"`
	}

	expectMap := make(map[string]interface{})
	expectMap["intKey"] = 123
	expectMap["strKey"] = "string"
	expectMap["floatKey"] = 234.5
	expectMap["sliceKey"] = []string{"a", "b"}
	expectMap["jsonKey"] = jK{
		Key: KV{
			K1: "v1",
			K2: "v2",
		},
	}

	for k, v := range expectMap {
		switch v.(type) {
		case int, int32, int64:
			require.Equal(t, v, scfg.GetInt(k))
		case string:
			require.Equal(t, v, scfg.GetString(k))
		case float64, float32:
			require.Equal(t, v, scfg.GetFloat64(k))
		case []string:
			require.Equal(t, v, scfg.GetStringSlice(k))
		default:
			actualVal := jK{}
			err := scfg.UnmarshalKey(k, &actualVal)
			require.Nil(t, err)
			require.Equal(t, v, actualVal)
		}

	}
}
