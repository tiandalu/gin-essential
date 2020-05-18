package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wcc4869/ginessential/common"
	"github.com/wcc4869/ginessential/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证 token，为空或者没有已 Bearer开头，报错
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足1",
				"data": tokenString,
			})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足2",
			})
			ctx.Abort()
			return
		}

		// 通过验证，获取 claims 中 userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 判断用户
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 用户存在,蒋用户信息写入上下文
		ctx.Set("user", user)

		ctx.Next()

	}
}
