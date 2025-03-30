package slogger

import (
	"context"
	"errors"
	"log/slog"
	"mylib/slogger/customhandle"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCustomLevelLogger(t *testing.T) {
	Convey("Given a log level, log type, and show source flag", t, func() {
		level := "Debug"
		logType := "TEXT"
		showSource := true
		ctxReadType := ContextReadTypeAllKeys

		Convey("When creating a new custom level logger", func() {
			logger := NewCustomLevelLogger(&Config{
				Level:       level,
				Format:      logType,
				ShowSource:  showSource,
				CtxReadType: ctxReadType,
			})

			Convey("Then the logger should not be nil", func() {
				So(logger, ShouldNotBeNil)
			})
		})
	})

	Convey("使用例子混杂", t, func() {
		// slog (默认日志级别为INFO)
		slog.Debug("slog: Hello, World!")
		slog.Info("slog: Hello, World!")
		slog.Warn("slog: Hello, World!")
		slog.Error("slog: Hello, World!")
		// slog with context
		ctx := context.WithValue(context.Background(), "trace_id", "123456")
		slog.DebugContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
		slog.InfoContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
		slog.WarnContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
		slog.ErrorContext(ctx, "slog: Hello, World!", "level", "LevelWarn")

		// slog with fields
		slog.With(slog.String("trace_id", "123123")).Info("slog: Hello, World!")
		slog.With(slog.String("trace_id", "123123"), slog.String("user_id", "123")).Info("slog: Hello, World!")

		// slog with JSON
		// level is DEBUG
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		logger.With(slog.String("trace_id", "123123"), slog.String("user_id", "123")).Debug("slog: Hello, World!")
		logger.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Debug("slog: Hello, World!")

		// set default logger
		slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "before")
		slog.SetDefault(logger)
		slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "after")

		// open source开启位置记录
		loggerWithSource := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
		slog.SetDefault(loggerWithSource)
		slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "after")

		// 添加自定义的json结构体作为字段
		slog.With(slog.Group("user", slog.String("name", "xiaoMing"), slog.Int("age", 18))).Warn("slog: with user info json", "gg", 18)

		// 自定义日志级别
		slog.SetDefault(NewCustomLevelLogger(&Config{
			Level:      "TRACE",
			Format:     LogFormatJSON,
			ShowSource: true,
		}))
		slog.DebugContext(ctx, "slog: Tracing 可以打debug（-4）", "aaa", "bbb")
		slog.Log(ctx, LevelTrace, "slog: Tracing 可以打Tracing（-2）")
		slog.InfoContext(ctx, "slog: Tracing 可以打Info（0）")
		slog.SetDefault(NewCustomLevelLogger(&Config{
			Level:      "FATAL",
			Format:     LogFormatTEXT,
			ShowSource: true,
		}))
		slog.DebugContext(ctx, "slog: Fatal 不可以打debug（-4）")
		slog.Log(ctx, LevelTrace, "slog: Fatal 不可以打Tracing（-2）")
		slog.InfoContext(ctx, "slog: Fatal 不可以打Info（0）")
		slog.Log(ctx, LevelFatal, "slog: Fatal 可以打Fatal（12）")

		// 使用不同的logger打印不同的逻辑
		logger1 := NewCustomLevelLogger(&Config{
			Level:       "INFO",
			Format:      LogFormatTEXT,
			ShowSource:  false,
			CtxReadType: ContextReadTypeAllKeys,
		})
		logger2 := NewCustomLevelLogger(&Config{
			Level:       "ERROR",
			Format:      LogFormatJSON,
			ShowSource:  true,
			CtxReadType: ContextReadTypeAllKeys,
		})

		logger1.Info("logger1: Hello, World!")
		logger2.Error("logger2: Hello, World!")

		// error的打印
		err := errors.New("slog: error message")
		logger1.With("error", err).Error("logger1: error message")
		logger2.With("error", err).Error("logger2: error message")

		// 补充context内容的打印
		myCtx := context.WithValue(context.Background(), "r_id", "r_123456")
		myCtx = context.WithValue(myCtx, "user", &Student{Name: "xiaoMing", Age: 18})
		logger1.ErrorContext(myCtx, "logger1: error message")
		logger2.ErrorContext(myCtx, "logger2: error message")
	})

	Convey("Given a log with prettyHandler", t, func() {
		opts := customhandle.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := customhandle.NewPrettyHandler(os.Stdout, opts)
		logger := slog.New(handler)
		logger.Debug(
			"executing database query",
			slog.String("query", "SELECT * FROM users"),
		)
		logger.Info("image upload successful", slog.String("image_id", "39ud88"))
		logger.Warn(
			"storage is 90% full",
			slog.String("available_space", "900.1 MB"),
		)
		logger.Error(
			"An error occurred while processing the request",
			slog.String("url", "https://example.com"),
		)
	})
}

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Addr *Addr
}

type Addr struct {
	Province string `json:"province"`
	City     string `json:"city"`
}

func TestGetKeyValues(t *testing.T) {

	Convey("Given a log level, log type, and show source flag", t, func() {
		ctx := context.WithValue(context.Background(), "trace_id", "123456")
		m := GetKeyValues(ctx)
		So(m, ShouldNotBeEmpty)
		So(m["trace_id"], ShouldEqual, "123456")
	})
}

// 打印context内容
func TestLogContextContent(t *testing.T) {
	Convey("should print successfully", t, func() {
		Convey("when CtxReadType=AllKeys", func() {
			sl := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatTEXT,
				ShowSource:  false,
				CtxReadType: ContextReadTypeAllKeys,
			})
			Convey("given context key is string, context value is string", func() {
				ctx := context.WithValue(context.Background(), "trace_id", "123456")
				sl.InfoContext(ctx, "slog: Hello, World!")
			})

			Convey("given context key is string, context value is struct", func() {
				ctx := context.WithValue(context.Background(), "user", &Student{Name: "xiaoMing", Age: 18})
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context key is struct, context value is struct", func() {
				ctx := context.WithValue(context.Background(), &Student{Name: "xiaoMing", Age: 18, Addr: &Addr{Province: "GuangDong", City: "ShenZhen"}}, &Student{Name: "xiaoMing", Age: 18, Addr: &Addr{Province: "GuangDong", City: "ShenZhen"}})
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context is only context.Background()", func() {
				ctx := context.Background()
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context is timeout ctx", func() {
				ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Hour)
				defer cancelFunc()
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context is cancel ctx", func() {
				ctx, cancelFunc := context.WithCancel(context.Background())
				defer cancelFunc()
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context is nil", func() {
				sl.InfoContext(nil, "slog: Hello, World!")
			})
		})

		Convey("when CtxReadType=Default", func() {
			sl := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatTEXT,
				ShowSource:  false,
				CtxReadKeys: []string{"trace_id"},
			})
			Convey("given context key is request_id, context value is 123456 [打印]", func() {
				ctx := context.WithValue(context.Background(), "request_id", "123456")
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
			Convey("given context key is trace_id, context value is 123456 [不打印]", func() {
				ctx := context.WithValue(context.Background(), "trace_id", "123456")
				sl.InfoContext(ctx, "slog: Hello, World!")
			})
		})

	})
	Convey("should panic", t, func() {
		Convey("given context key is map, context value is string", func() {
			m := map[string]interface{}{
				"trace_id": "123456",
			}
			// 本身这句话就会panic
			So(func() {
				context.WithValue(context.Background(), m, "123456")
			}, ShouldPanicWith, "key is not comparable")
		})
		Convey("given context key is slice, context value is struct", func() {
			s := []string{"trace_id", "123456"}
			// 本身这句话就会panic
			So(func() {
				context.WithValue(context.Background(), s, &Student{Name: "xiaoMing", Age: 18})
			}, ShouldPanicWith, "key is not comparable")
		})
	})
}

var CtxReqID = struct{}{}

func TestAddCtxReadKeyStruct(t *testing.T) {
	Convey("should 打印这个ctx内部的struct{}{} 并以 key=request_id 的方式打印", t, func() {
		Convey("given context key is struct, context value is 123456", func() {
			lg := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatJSON,
				ShowSource:  true,
				CtxReadType: ContextReadTypeDefault,
			})
			AddCtxReadKeyStruct(lg, CtxReqID, "request_id")

			ctx := context.WithValue(context.Background(), CtxReqID, "123456")

			lg.InfoContext(ctx, "slog: Hello, World!")
		})
	})
	Convey("should 不打印request_id", t, func() {
		Convey("given context key is struct, context value is 123456 , 正常没有 AddCtxReadKeyStruct", func() {
			lg := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatJSON,
				ShowSource:  true,
				CtxReadType: ContextReadTypeDefault,
			})
			// addCtxReadKey(lg, CtxReqID, "request_id")

			ctx := context.WithValue(context.Background(), CtxReqID, "123456")

			lg.InfoContext(ctx, "slog: Hello, World!")
		})
	})
}

func TestAddCtxReadKeyString(t *testing.T) {
	Convey("should 打印这个ctx内部的string 并以 key=request_id 的方式打印", t, func() {
		Convey("given context key is struct, context value is 123456", func() {
			lg := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatJSON,
				ShowSource:  true,
				CtxReadType: ContextReadTypeDefault,
			})
			AddCtxReadKeyString(lg, "x_id")

			ctx := context.WithValue(context.Background(), "x_id", "bbbb")
			ctx = context.WithValue(ctx, "request_id", "123456")

			lg.InfoContext(ctx, "slog: Hello, World!")
		})
	})
	Convey("should [并发读写场景]打印这个ctx内部的string 并以 key=request_id 的方式打印", t, func() {
		Convey("given context key is struct, context value is 123456", func() {
			lg := NewCustomLevelLogger(&Config{
				Level:       "INFO",
				Format:      LogFormatJSON,
				ShowSource:  true,
				CtxReadType: ContextReadTypeDefault,
			})

			ctx := context.WithValue(context.Background(), "x_id", "bbbb")
			ctx = context.WithValue(ctx, "request_id", "123456")

			wg := sync.WaitGroup{}
			wg.Add(6)
			go func() {
				for i := 0; i < 10000; i++ {
					AddCtxReadKeyString(lg, "x_id")
				}
				wg.Done()
			}()
			go func() {
				for i := 0; i < 10000; i++ {
					AddCtxReadKeyString(lg, "x_id")
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 10000; i++ {
					AddCtxReadKeyString(lg, "x_id")
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 10000; i++ {
					lg.InfoContext(ctx, "slog: Hello, World!")
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 10000; i++ {
					lg.InfoContext(ctx, "slog: Hello, World!")
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 10000; i++ {
					lg.InfoContext(ctx, "slog: Hello, World!")
				}
				wg.Done()
			}()

			wg.Wait()
		})
	})

}
