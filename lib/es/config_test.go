package es

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func GetTestCfg() *Config {
	// 开发环境es
	// kibana：http://10.160.12.33:5601
	return &Config{
		Enable:   true,
		URL:      "http://10.160.12.33:9200",
		Username: "wps",
		Password: "jW8Bk2jezx9Ssp2",
	}
	//return &Config{
	//	Enable: true,
	//	URL:    "http://10.13.148.26:9200",
	//	//Username: "wps",
	//	//Password: "jW8Bk2jezx9Ssp2",
	//}
}

func TestConfig_BuildES7Client(t *testing.T) {
	Convey("Given a valid Elasticsearch configuration", t, func() {
		cli, err := GetTestCfg().BuildES7Client()
		So(err, ShouldBeNil)
		So(cli, ShouldNotBeNil)
	})
}

func TestConfig_BuildES8Client(t *testing.T) {
	Convey("Given a valid Elasticsearch configuration", t, func() {
		cli, err := GetTestCfg().BuildES8Client()
		So(err, ShouldBeNil)
		So(cli, ShouldNotBeNil)
	})
}
