package account

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Msg string `json:"msg"`
}

func Login(c *gin.Context, req *LoginRequest) (*LoginResponse, error) {
	if len(req.Name) == 0 {
		return nil, errors.New("param name is empty")
	}
	return &LoginResponse{Msg: "ok"}, nil
}
