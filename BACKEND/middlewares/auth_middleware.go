package middlewares

import (
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 中间件是为了过滤路由的一种机制，在到达处理函数之前， 拦截请求并进行相应处理(比如: 权限验证, 数据过滤等), 这个可以类比为前置拦截器或前置过滤器
func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorization"})
			return
		}
		username, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
		ctx.Set("username", username)
		ctx.Next()
	}
}