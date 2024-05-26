package swlimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

const MaxVisitTimes = 10

type visit struct {
	// 最后一次请求时间
	lastVisit time.Time

	// 对应Time窗口内的访问次数
	visitTimes int
}

// 全局map来存储所有IP的请求信息
var visitMap = sync.Map{}

func LimitFunc(ctx *gin.Context) {
	ip := ctx.Request.RemoteAddr
	v, ok := visitMap.Load(ip)
	if !ok {
		// 首次请求，则初始化visit
		visitMap.Store(ip, &visit{lastVisit: time.Now(), visitTimes: 1})
		ctx.Next()
		return
	}
	// 若该IP非首次请求，且距离上次请求时间超过Time窗口，则重设visitTimes
	visitObj := v.(*visit)
	fmt.Println("ip["+ip+"]当前的visitObj是：", visitObj.lastVisit, visitObj.visitTimes)
	fmt.Println("since：", time.Since(visitObj.lastVisit))

	if time.Since(visitObj.lastVisit) > time.Minute {
		// 每分钟重置
		visitObj.visitTimes = 1
	} else if visitObj.visitTimes > MaxVisitTimes {
		// 若本次请求距离上次请求时间在Time窗口内，且该IP在此时间内的访问次数超过上限，则返回错误
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"msg": "too many requests"})
		return
	} else {
		visitObj.visitTimes++
	}
	visitObj.lastVisit = time.Now()

	ctx.Next()
}
