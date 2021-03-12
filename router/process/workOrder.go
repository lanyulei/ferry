/*
  @Author : lanyulei
*/

package process

import (
	"ferry/apis/process"
	"ferry/middleware"
	jwt "ferry/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

func RegisterWorkOrderRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	workOrderRouter := v1.Group("/work-order").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		workOrderRouter.GET("/process-structure", process.ProcessStructure)
		workOrderRouter.POST("/create", process.CreateWorkOrder)
		workOrderRouter.GET("/list", process.WorkOrderList)
		workOrderRouter.POST("/handle", process.ProcessWorkOrder)
		workOrderRouter.GET("/unity", process.UnityWorkOrder)
		workOrderRouter.POST("/inversion", process.InversionWorkOrder)
		workOrderRouter.GET("/urge", process.UrgeWorkOrder)
		workOrderRouter.PUT("/active-order/:id", process.ActiveOrder)
		workOrderRouter.DELETE("/delete/:id", process.DeleteWorkOrder)
		workOrderRouter.POST("/reopen/:id", process.ReopenWorkOrder)
	}
}
