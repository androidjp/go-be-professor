package apollocli

import (
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"mylib/scfg/mgr"

	"github.com/stretchr/testify/assert"
)

var (
	TestConfigServerURL         = "http://10.160.12.32:8080" // 以真实配置服务作为测试用例，若测试服务器下线则需要更换
	TestAppID                   = "apollo-client-test"
	TestCluster                 = "default"
	TestSubscribeNamespaceNames = []string{
		"application",
		"test-ns",
		"test-json.json",
	}
	TestDefaultNamespace = "application"
)

type Case struct {
	key      string
	expected interface{}
}

func TestApolloClient(t *testing.T) {
	cli := NewApolloClientProxy(&ApolloConfig{
		AppID:                   TestAppID,
		Cluster:                 TestCluster,
		SubscribeNamespaceNames: TestSubscribeNamespaceNames,
		DefaultNamespace:        TestDefaultNamespace,
		IP:                      TestConfigServerURL,
	},
		WithSysLoggerMgr(mgr.NewSysLogMgr(slog.Default(), true)),
	)

	assert.Nil(t, cli.Start())

	// 按顺序来，依赖配置中心配置的测试值
	testCases := map[string]Case{
		"int": {
			key:      "int",
			expected: 1,
		},
		"int32": {
			key:      "int32",
			expected: int32(2),
		},
		"int64": {
			key:      "int64",
			expected: int64(3),
		},
		"uint": {
			key:      "uint",
			expected: uint(4),
		},
		"uint32": {
			key:      "uint32",
			expected: uint32(5),
		},
		"uint64": {
			key:      "uint64",
			expected: uint64(6),
		},
		"float64": {
			key:      "float64",
			expected: 7.2,
		},
		"string": {
			key:      "string",
			expected: "8",
		},
		"bool": {
			key:      "bool",
			expected: false,
		},
		"intslice": {
			key:      "intslice",
			expected: []int{1, 2, 3},
		},
		"stringslice": {
			key:      "stringslice",
			expected: []string{"1", "2", "3"},
		},
		"stringmapstring": {
			key: "stringmapstring",
			expected: map[string]string{
				"a": "a1",
				"b": "b1",
			},
		},
		"stringmapstringslice": {
			key: "stringmapstringslice",
			expected: map[string][]string{
				"a": {"a1"},
			},
		},
	}

	// 按顺序来，依赖配置中心配置的测试值
	assert.Equal(t, testCases["int"].expected.(int), cli.GetInt(testCases["int"].key, 2))
	assert.Equal(t, testCases["int32"].expected.(int32), cli.GetInt32(testCases["int32"].key, 2))
	assert.Equal(t, testCases["int64"].expected.(int64), cli.GetInt64(testCases["int64"].key, 2))
	assert.Equal(t, testCases["uint"].expected.(uint), cli.GetUInt(testCases["uint"].key, 2))
	assert.Equal(t, testCases["uint32"].expected.(uint32), cli.GetUInt32(testCases["uint32"].key, 2))
	assert.Equal(t, testCases["uint64"].expected.(uint64), cli.GetUInt64(testCases["uint64"].key, 2))
	assert.Equal(t, testCases["float64"].expected.(float64), cli.GetFloat64(testCases["float64"].key, 2))
	assert.Equal(t, testCases["string"].expected.(string), cli.GetString(testCases["string"].key, "2"))
	assert.Equal(t, testCases["bool"].expected.(bool), cli.GetBool(testCases["bool"].key, true))
	assert.Equal(t, testCases["intslice"].expected.([]int), cli.GetIntSlice(testCases["intslice"].key, []int{1}))
	assert.Equal(t, testCases["stringslice"].expected.([]string), cli.GetStringSlice(testCases["stringslice"].key, []string{"1"}))
	assert.Equal(t, testCases["stringmapstring"].expected.(map[string]string), cli.GetStringMapString(testCases["stringmapstring"].key, map[string]string{
		"a": "a1",
	}))

	assert.Equal(t, testCases["stringmapstringslice"].expected.(map[string][]string), cli.GetStringMapStringSlice(testCases["stringmapstringslice"].key, map[string][]string{
		"b": {"a1"},
	}))
	assert.Equal(t, "test-ns", cli.GetString("test", "adda", WithNamespace("test-ns")))

	expected := map[string]string{
		"json": "json namespace!",
	}

	var actual map[string]string
	assert.Nil(t, json.Unmarshal([]byte(cli.GetContent(WithNamespace("test-json.json"))), &actual))
	assert.Equal(t, expected, actual)
}

func TestJSONUnmarshal(t *testing.T) {
	cli := NewApolloClientProxy(&ApolloConfig{
		AppID:                   TestAppID,
		Cluster:                 TestCluster,
		SubscribeNamespaceNames: TestSubscribeNamespaceNames,
		DefaultNamespace:        TestDefaultNamespace,
		IP:                      TestConfigServerURL,
	})
	assert.Nil(t, cli.Start())

	cli.SettingSysLogMgr(mgr.NewSysLogMgr(slog.Default(), true))

	type TargetObject struct {
		A     string        `json:"a"`
		Time  time.Duration `json:"time"`
		Slice []int         `json:"slice"`
	}

	var target TargetObject
	key := "jsonobj"

	expectObject := TargetObject{
		A:     "a",
		Time:  time.Second,
		Slice: []int{1, 2, 3},
	}

	assert.Nil(t, cli.UnmarshalKey(key, &target))
	assert.Equal(t, expectObject, target)
}

func TestApolloClientDefaultNamespaceEmptyCase(t *testing.T) {
	cli := NewApolloClientProxy(&ApolloConfig{
		AppID:                   TestAppID,
		Cluster:                 TestCluster,
		SubscribeNamespaceNames: TestSubscribeNamespaceNames,
		DefaultNamespace:        "",
		IP:                      TestConfigServerURL,
	})
	assert.Equal(t, "application.properties", cli.defaultNamespace)
}

func TestApolloClientValueEmptyCase(t *testing.T) {
	cli := NewApolloClientProxy(&ApolloConfig{
		AppID:                   TestAppID,
		Cluster:                 TestCluster,
		SubscribeNamespaceNames: TestSubscribeNamespaceNames,
		DefaultNamespace:        TestDefaultNamespace,
		IP:                      TestConfigServerURL,
	})
	assert.Nil(t, cli.Start())

	defaultValue := "default value"
	actual := cli.GetString("not_exists_key", defaultValue)
	assert.Equal(t, defaultValue, actual)
}

// 测试获取不到远程配置去默认值的用例
func TestApolloClientDefaultValues(t *testing.T) {
	cli := NewApolloClientProxy(&ApolloConfig{
		AppID:                   TestAppID,
		Cluster:                 TestCluster,
		SubscribeNamespaceNames: TestSubscribeNamespaceNames,
		DefaultNamespace:        "Invalid Namespace",
		IP:                      TestConfigServerURL,
	},
		WithSysLoggerMgr(mgr.NewSysLogMgr(slog.Default(), true)),
	)

	assert.Nil(t, cli.Start())

	testCases := map[string]Case{
		"int": {
			key:      "int",
			expected: 1,
		},
		"int32": {
			key:      "int32",
			expected: int32(2),
		},
		"int64": {
			key:      "int64",
			expected: int64(3),
		},
		"uint": {
			key:      "uint",
			expected: uint(4),
		},
		"uint32": {
			key:      "uint32",
			expected: uint32(5),
		},
		"uint64": {
			key:      "uint64",
			expected: uint64(6),
		},
		"float64": {
			key:      "float64",
			expected: 7.2,
		},
		"string": {
			key:      "string",
			expected: "8",
		},
		"bool": {
			key:      "bool",
			expected: false,
		},
		"intslice": {
			key:      "intslice",
			expected: []int{1, 2, 3},
		},
		"stringslice": {
			key:      "stringslice",
			expected: []string{"1", "2", "3"},
		},
		"stringmapstring": {
			key: "stringmapstring",
			expected: map[string]string{
				"a": "a1",
				"b": "b1",
			},
		},
		"stringmapstringslice": {
			key: "stringmapstringslice",
			expected: map[string][]string{
				"a": {"a1"},
			},
		},
	}

	// 按顺序来，依赖配置中心配置的测试值
	assert.Equal(t, testCases["int"].expected.(int), cli.GetInt(testCases["int"].key, testCases["int"].expected.(int)))
	assert.Equal(t, testCases["int32"].expected.(int32), cli.GetInt32(testCases["int32"].key, testCases["int32"].expected.(int32)))
	assert.Equal(t, testCases["int64"].expected.(int64), cli.GetInt64(testCases["int64"].key, testCases["int64"].expected.(int64)))
	assert.Equal(t, testCases["uint"].expected.(uint), cli.GetUInt(testCases["uint"].key, testCases["uint"].expected.(uint)))
	assert.Equal(t, testCases["uint32"].expected.(uint32), cli.GetUInt32(testCases["uint32"].key, testCases["uint32"].expected.(uint32)))
	assert.Equal(t, testCases["uint64"].expected.(uint64), cli.GetUInt64(testCases["uint64"].key, testCases["uint64"].expected.(uint64)))
	assert.Equal(t, testCases["float64"].expected.(float64), cli.GetFloat64(testCases["float64"].key, testCases["float64"].expected.(float64)))
	assert.Equal(t, testCases["string"].expected.(string), cli.GetString(testCases["string"].key, testCases["string"].expected.(string)))
	assert.Equal(t, testCases["bool"].expected.(bool), cli.GetBool(testCases["bool"].key, testCases["bool"].expected.(bool)))
	assert.Equal(t, testCases["intslice"].expected.([]int), cli.GetIntSlice(testCases["intslice"].key, testCases["intslice"].expected.([]int)))
	assert.Equal(t, testCases["stringslice"].expected.([]string), cli.GetStringSlice(testCases["stringslice"].key, testCases["stringslice"].expected.([]string)))
	assert.Equal(t, testCases["stringmapstring"].expected.(map[string]string), cli.GetStringMapString(testCases["stringmapstring"].key, testCases["stringmapstring"].expected.(map[string]string)))

	assert.Equal(t, testCases["stringmapstringslice"].expected.(map[string][]string), cli.GetStringMapStringSlice(testCases["stringmapstringslice"].key, testCases["stringmapstringslice"].expected.(map[string][]string)))
}
