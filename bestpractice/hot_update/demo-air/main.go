package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong123",
		})
	})

	if err := r.Run("localhost:8080"); err != nil {
		fmt.Println(err)
	}
}
