package order

import (
	"demo-gin/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListResponse struct {
	Items []string `json:"items"`
}

func List(c *gin.Context, req *helper.Empty) (*ListResponse, error) {
	c.JSON(http.StatusOK, gin.H{"msg": "xxxx"})
	return nil, nil
}
