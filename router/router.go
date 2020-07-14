package router

import (
	"ferry/apis/monitor"
	"ferry/apis/system"
	"ferry/handler"
	"ferry/pkg/jwtauth"
	jwt "ferry/pkg/jwtauth"
	"ferry/router/process"
	systemRouter "ferry/router/system"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")

	sysBaseRouter(g)
	// 静态文件
	sysStaticFileRouter(g)

	// swagger；注意：生产环境可以注释掉
	sysSwaggerRouter(g)

	// 无需认证
	sysNoCheckRoleRouter(g)
	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)

	return g
}

func sysBaseRouter(r *gin.RouterGroup) {
	r.GET("/", system.HelloWorld)
	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	r.Static("/static", "./static")
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func sysNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")

	v1.GET("/monitor/server", monitor.ServerInfo)
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwtauth.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	v1 := r.Group("/api/v1")

	// 系统管理
	systemRouter.RegisterPageRouter(v1, authMiddleware)
	systemRouter.RegisterBaseRouter(v1, authMiddleware)
	systemRouter.RegisterDeptRouter(v1, authMiddleware)
	systemRouter.RegisterSysUserRouter(v1, authMiddleware)
	systemRouter.RegisterRoleRouter(v1, authMiddleware)
	systemRouter.RegisterUserCenterRouter(v1, authMiddleware)
	systemRouter.RegisterPostRouter(v1, authMiddleware)
	systemRouter.RegisterMenuRouter(v1, authMiddleware)
	systemRouter.RegisterLoginLogRouter(v1, authMiddleware)

	// 流程中心
	process.RegisterClassifyRouter(v1, authMiddleware)
	process.RegisterProcessRouter(v1, authMiddleware)
	process.RegisterTaskRouter(v1, authMiddleware)
	process.RegisterTplRouter(v1, authMiddleware)
	process.RegisterWorkOrderRouter(v1, authMiddleware)
}
