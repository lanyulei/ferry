package system

import (
	log2 "ferry/apis/log"
	"ferry/apis/monitor"
	"ferry/apis/public"
	"ferry/apis/system"
	_ "ferry/docs"
	"ferry/handler"
	"ferry/middleware"
	jwt "ferry/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

func SysBaseRouter(r *gin.RouterGroup) {
	//r.GET("/", system.HelloWorld)
	r.GET("/info", handler.Ping)
}

func SysNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")

	v1.GET("/monitor/server", monitor.ServerInfo)
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/settings", system.GetSettingsInfo)

	registerPublicRouter(v1)
}

func RegisterBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/getinfo", system.GetInfo)
		v1auth.GET("/menurole", system.GetMenuRole)
		v1auth.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
		v1auth.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleSelect)
		v1auth.POST("/logout", handler.LogOut)
		v1auth.GET("/menuids", system.GetMenuIDS)
	}
}

func RegisterPageRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/deptList", system.GetDeptList)
		v1auth.GET("/ordinaryDeptList", system.GetOrdinaryDeptList)
		v1auth.GET("/deptTree", system.GetDeptTree)
		v1auth.GET("/sysUserList", system.GetSysUserList)
		v1auth.GET("/rolelist", system.GetRoleList)
		v1auth.GET("/postlist", system.GetPostList)
		v1auth.GET("/menulist", system.GetMenuList)
		v1auth.GET("/loginloglist", log2.GetLoginLogList)
	}
}

func RegisterUserCenterRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	user := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		user.GET("/profile", system.GetSysUserProfile)
		user.POST("/avatar", system.InsetSysUserAvatar)
		user.PUT("/pwd", system.SysUserUpdatePwd)
	}
}

func RegisterLoginLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	loginlog := v1.Group("/loginlog").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		loginlog.GET("/:infoId", log2.GetLoginLog)
		loginlog.POST("", log2.InsertLoginLog)
		loginlog.PUT("", log2.UpdateLoginLog)
		loginlog.DELETE("/:infoId", log2.DeleteLoginLog)
		loginlog.DELETE("", log2.CleanLoginLog)
	}
}

func RegisterPostRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	post := v1.Group("/post").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		post.GET("/:postId", system.GetPost)
		post.POST("", system.InsertPost)
		post.PUT("", system.UpdatePost)
		post.DELETE("/:postId", system.DeletePost)
	}
}

func RegisterMenuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	menu := v1.Group("/menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		menu.GET("/:id", system.GetMenu)
		menu.POST("", system.InsertMenu)
		menu.PUT("", system.UpdateMenu)
		menu.DELETE("/:id", system.DeleteMenu)
	}
}

func RegisterRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	role := v1.Group("/role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		role.GET("/:roleId", system.GetRole)
		role.POST("", system.InsertRole)
		role.PUT("", system.UpdateRole)
		role.DELETE("/:roleId", system.DeleteRole)
	}
}

func RegisterSysUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	sysuser := v1.Group("/sysUser").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		sysuser.GET("/:userId", system.GetSysUser)
		sysuser.GET("/", system.GetSysUserInit)
		sysuser.POST("", system.InsertSysUser)
		sysuser.PUT("", system.UpdateSysUser)
		sysuser.DELETE("/:userId", system.DeleteSysUser)
	}
}

func RegisterDeptRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dept := v1.Group("/dept").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		dept.GET("/:deptId", system.GetDept)
		dept.POST("", system.InsertDept)
		dept.PUT("", system.UpdateDept)
		dept.DELETE("/:id", system.DeleteDept)
	}
}

func RegisterSysSettingRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	setting := v1.Group("/settings").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		setting.POST("", system.SetSettingsInfo)
	}
}

func registerPublicRouter(v1 *gin.RouterGroup) {
	p := v1.Group("/public")
	{
		p.POST("/uploadFile", public.UploadFile)
	}
}
