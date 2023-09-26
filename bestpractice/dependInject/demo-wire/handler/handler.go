package handler

import (
	"demo-wire/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

type PostHandler struct {
	serv service.IPostService // 成员对象，是一个接口类型
}

func (h *PostHandler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/post/:id", h.GetPostById)
}
func (h *PostHandler) GetPostById(ctx *gin.Context) {
	content := h.serv.GetPostById(ctx, ctx.Param("id"))
	ctx.String(http.StatusOK, content)
}

// 这可以看做是一个wire提供者
func NewPostHandler(serv service.IPostService) *PostHandler {
	return &PostHandler{serv: serv}
}

// v2版本
var PostSet = wire.NewSet(NewPostHandler, service.NewPostService)
