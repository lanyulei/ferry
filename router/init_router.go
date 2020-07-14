package router

import (
	"ferry/handler"
	"ferry/middleware"
	_ "ferry/pkg/jwtauth"
	"ferry/tools"
	config2 "ferry/tools/config"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	if config2.ApplicationConfig.IsHttps {
		r.Use(handler.TlsHandler())
	}
	middleware.InitMiddleware(r)
	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	return r
}
