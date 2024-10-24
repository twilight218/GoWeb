package controllers

import (
	"exchangeapp/global"
	"exchangeapp/model"
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ctx *gin.Context 表示一个请求的上下文，它包含了请求和响应的所有信息。
// ctx 是 Gin 框架中的核心对象，通过它可以访问请求数据（如参数、头信息、体数据等），也可以发送响应数据。
// 注册响应函数，需要一个上下文参数
func Register(ctx *gin.Context) {
	var user model.User
	// ShouldBindJSON 方法用于尝试将请求体中的JSON数据绑定到指定的Go结构体上
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}
	// 密码加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}
	user.Password = hashedPwd
	// 使用jwt生成用户的token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}
	// AutoMigrate没有表会创建与类一致的表，有表了会根据模型更新表结构
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {"error": err.Error()})
		return
	}

	// 创建记录
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H {
		"token": token,
	})
}

func Login(ctx *gin.Context) {
	// 结构体中，成员后面使用反引号定义描述性信息tag，用于序列化和反序列化
	// 用于接收登录信息
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	var user model.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H {"error": "illegal username"})
		return
	}
	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H {"error": "wrong password"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}