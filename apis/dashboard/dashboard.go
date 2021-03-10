package dashboard

import (
	"ferry/pkg/service"
	"ferry/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func InitData(c *gin.Context) {
	var (
		err            error
		workOrderCount map[string]int // 工单数量统计
	)

	workOrderCount, err = service.WorkOrderCount(c)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	app.OK(c, map[string]interface{}{
		"workOrderCount": workOrderCount,
	}, "")
}
