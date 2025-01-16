package main

import (
	"demo-gin/handler/account"
	"demo-gin/handler/order"
	"demo-gin/helper"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	{
		helper.POST(g, "acc/login", account.Login)
	}
	{
		helper.POST(g, "order/list", order.List)
		helper.POST(g, "order/get", order.Get)
	}

	g.Run(":8000")
}
