package redis

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig_Build(t *testing.T) {
	Convey("should 返回 开发环境 127.0.0.1:6379 单例redis [前提：本地启动了redis]", t, func() {
		Convey("given DefaultConfig, config使用sample方式Build", func() {
			cfg := DefaultConfig()
			redisCli, err := cfg.Build()
			So(err, ShouldBeNil)
			defer redisCli.Close()
			So(err, ShouldBeNil)
			set := redisCli.Set(context.Background(), "offcieclouddocsrv_key1", "val1", time.Second)
			So(set.Err(), ShouldBeNil)
			get := redisCli.Get(context.Background(), "offcieclouddocsrv_key1")
			So(get.Err(), ShouldBeNil)
			So(get.Val(), ShouldEqual, "val1")
		})
	})

	Convey("should 返回 返回nil", t, func() {
		Convey("given config.Enable = false", func() {
			cfg := DefaultConfig()
			cfg.Enable = false
			redisCli, err := cfg.Build()
			So(err, ShouldBeNil)
			So(redisCli, ShouldBeNil)
		})
	})
}

func TestConfig_BuildSentinelClient(t *testing.T) {

}
