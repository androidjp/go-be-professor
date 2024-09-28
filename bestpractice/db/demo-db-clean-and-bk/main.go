package main

import (
	"demodbclient/model"
	"demodbclient/repo"
	"demodbclient/task"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化repo
	studentRepo := repo.NewRepo()
	defer studentRepo.Stop()

	// 定时任务
	t := task.NewCleanTask(studentRepo)
	t.Run()

	// 模拟本地TiDB场景
	r := gin.Default()
	r.GET("/insert", func(c *gin.Context) {
		// 插入数据
		now := time.Now()
		err := studentRepo.InsertStudent(&model.Student{
			Name:  "Jasper",
			Age:   20,
			Ctime: now.Unix(),
			Mtime: now.Unix(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})
	r.GET("/batch_read", func(c *gin.Context) {
		// 查询数据
		stus, err := studentRepo.BatchRead()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": stus})
	})

	_ = r.Run("127.0.0.1:8080")
}
