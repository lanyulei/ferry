package process

import (
	"ferry/apis/process"
	"ferry/middleware"
	jwt "ferry/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RegisterTaskRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	taskRouter := v1.Group("/task").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		taskRouter.GET("", process.TaskList)
		taskRouter.GET("/details", process.TaskDetails)
		taskRouter.POST("", process.CreateTask)
		taskRouter.PUT("", process.UpdateTask)
		taskRouter.DELETE("", process.DeleteTask)
	}
}
