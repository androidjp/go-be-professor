package main

import (
	"demo-gorm-gen/dal"
	"demo-gorm-gen/dal/model"
	"demo-gorm-gen/dal/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	gin.ForceConsoleColor()

	dal.GenerateRepo("./dal/query")

	// 针对user的CRUD

	userG := r.Group("/api/v1/users")
	{
		// 批量查询用户
		userG.GET("", func(c *gin.Context) {
			var q = query.Use(dal.DB.Debug())
			ud := q.User.WithContext(c)

			if uid := c.Param("id"); len(uid) != 0 {
				id, _ := strconv.Atoi(uid)
				result, err := ud.FindByID(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"msg": err.Error(),
					})
					return
				}
				c.JSON(http.StatusOK, map[string]interface{}{
					"msg":  "success",
					"data": result,
				})
				return
			}

			offsetStr, _ := c.GetQuery("offset")
			limitStr, _ := c.GetQuery("limit")
			offset, _ := strconv.Atoi(offsetStr)

			limit, _ := strconv.Atoi(limitStr)
			if limit == 0 {
				limit = 10
			}

			result, count, err := ud.FindByPage(offset, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"msg":   "success",
				"data":  result,
				"count": count,
			})
		})

		// 查询具体用户
		userG.GET("/:id", func(c *gin.Context) {
			var q = query.Use(dal.DB.Debug())
			ud := q.User.WithContext(c)

			uid := c.Param("id")
			id, _ := strconv.Atoi(uid)
			result, err := ud.FindByID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"msg":  "success",
				"data": result,
			})
		})

		// 创建用户
		userG.POST("", func(c *gin.Context) {
			var q = query.Use(dal.DB.Debug())
			ud := q.User.WithContext(c)

			userData := &model.User{}
			if err := c.ShouldBindJSON(userData); err != nil {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}

			if err := ud.Create(userData); err != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, map[string]interface{}{
				"msg": "success",
			})
		})

		// 更新
		userG.PUT("", func(c *gin.Context) {
			var q = query.Use(dal.DB.Debug())
			ud := q.User.WithContext(c)

			userData := &model.User{}
			if err := c.ShouldBindJSON(userData); err != nil {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}
			if userData.ID == 0 {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"msg": "user id is invalid",
				})
				return
			}

			//  更新操作
			if err := ud.Save(userData); err != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, map[string]interface{}{
				"msg": "success",
			})
		})

		// 删除
		userG.DELETE("/:id", func(c *gin.Context) {
			var q = query.Use(dal.DB.Debug())
			ud := q.User.WithContext(c)

			uid := c.Param("id")
			id, _ := strconv.Atoi(uid)
			result, err := ud.Delete(&model.User{ID: uint(id)})
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"msg":  "success",
				"data": result,
			})
		})

	}

	r.Run(":8222")
}
