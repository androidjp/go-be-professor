package bootstrap

import (
	"context"
	"log/slog"
	"testing"

	"mylib/scfg"

	"github.com/smartystreets/goconvey/convey"
)

func TestInitGlobalCfgAndLog(t *testing.T) {
	convey.Convey("Given app-with-slogger.yaml", t, func() {
		convey.Convey("When InitGlobalCfgAndLog is called", func() {
			InitGlobalCfgAndLog("./../mock/cfg/app-with-slogger.yaml")

			convey.Convey("Then slog.Enabled(slog.LevelDebug) is false, slog.Enabled(slog.LevelInfo) is false, slog.Enabled(slog.LevelWarn) is true", func() {
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelDebug), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelInfo), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelWarn), convey.ShouldBeTrue)
				ctxWithReqID := context.WithValue(context.Background(), "request_id", "123456")
				ctxWithReqID = context.WithValue(ctxWithReqID, "trace_id", "123456")
				slog.With(slog.Group("user", "name", "Mike", "age", 18)).WarnContext(ctxWithReqID, "test with request_id and trace_id")
			})
		})
	})

	convey.Convey("Given app-empty.yaml", t, func() {
		convey.Convey("When InitGlobalCfgAndLog is called", func() {
			InitGlobalCfgAndLog("./../mock/cfg/app-empty.yaml")

			convey.Convey("Then slog.Enabled(slog.LevelDebug) is false, slog.Enabled(slog.LevelInfo) is true", func() {
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelDebug), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelInfo), convey.ShouldBeTrue)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelWarn), convey.ShouldBeTrue)
				ctxWithReqID := context.WithValue(context.Background(), "request_id", "123456")
				slog.WarnContext(ctxWithReqID, "test with request_id")
			})
		})
	})
}

func TestInitGlobalSlogBySCfg(t *testing.T) {
	convey.Convey("Given app-with-slogger.yaml", t, func() {
		scfg.ParseConfig("./../mock/cfg/app-with-slogger.yaml", slog.Default())

		convey.Convey("When InitGlobalSlogBySCfg is called", func() {
			initGlobalSlogBySCfg()

			convey.Convey("Then slog.Enabled(slog.LevelDebug) is false, slog.Enabled(slog.LevelInfo) is false, slog.Enabled(slog.LevelWarn) is true", func() {
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelDebug), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelInfo), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelWarn), convey.ShouldBeTrue)
			})
		})
	})

	convey.Convey("Given app-empty.yaml", t, func() {
		scfg.ParseConfig("./../mock/cfg/app-empty.yaml", slog.Default())

		convey.Convey("When InitGlobalSlogBySCfg is called", func() {
			initGlobalSlogBySCfg()

			convey.Convey("Then slog.Enabled(slog.LevelDebug) is false, slog.Enabled(slog.LevelInfo) is true", func() {
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelDebug), convey.ShouldBeFalse)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelInfo), convey.ShouldBeTrue)
				convey.So(slog.Default().Enabled(context.Background(), slog.LevelWarn), convey.ShouldBeTrue)
			})
		})
	})
}
