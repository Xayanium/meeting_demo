package middlewares

import (
	"github.com/gin-gonic/gin"
	"meeting_demo/utils"
	"net/http"
)

// Auth 鉴权中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization") // 解析请求头中的Authorization字段

		userAuth, err := utils.AnalyseToken(tokenStr) // 解析token字符串，得到用户声明
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, &gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authorization",
			})
			return
		}

		c.Set("user_token", userAuth) // 将解析出的用户声明存储在上下文中，以便后续的处理函数可以访问这些数据
		c.Next()                      // 继续执行下一个中间件或路由处理器
	}
}
