package main

import (
	"fmt"
	"metricdemo/metrics"
	"metricdemo/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"
const ProjectAccessLogFile = "./logs/access.log"

func main() {
	r := gin.Default()

	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", "jp-test", "dev")),
		logger.WithTimeLayout(CSTLayout),
		logger.WithFileP(ProjectAccessLogFile),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()

	// 初始化promethous 客户端
	r.Use(func(c *gin.Context) {
		ts := time.Now()

		defer metrics.RecordHandler(accessLogger)(&metrics.MetricsMessage{
			ProjectName:  "hello",
			Env:          "dev",
			TraceID:      ts.String(),
			HOST:         c.Request.Host,
			Path:         c.Request.RequestURI,
			Method:       c.Request.Method,
			HTTPCode:     c.Writer.Status(),
			BusinessCode: 200,
			CostSeconds:  time.Since(ts).Seconds(),
			IsSuccess:    !c.IsAborted() && (c.Writer.Status() == http.StatusOK),
		})

		c.Next()
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "golang ~")
	})

	http.ListenAndServe(":8080", r)
}
