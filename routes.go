package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wcc4869/ginessential/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	return r
}
