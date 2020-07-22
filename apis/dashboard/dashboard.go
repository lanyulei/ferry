package dashboard

import (
	"ferry/global/orm"
	"ferry/models/process"
	"ferry/models/system"
	"ferry/pkg/pagination"
	"ferry/pkg/service"
	"ferry/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func InitData(c *gin.Context) {
	var (
		err        error
		panelGroup struct {
			UserTotalCount      int `json:"user_total_count"`
			WorkOrderTotalCount int `json:"work_order_total_count"`
			UpcomingTotalCount  int `json:"upcoming_total_count"`
			MyUpcomingCount     int `json:"my_upcoming_count"`
		}
		result interface{}
	)

	// 查询用户总数
	err = orm.Eloquent.Model(&system.SysUser{}).Count(&panelGroup.UserTotalCount).Error
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 查询工单总数
	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).Count(&panelGroup.WorkOrderTotalCount).Error
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 查询待办总数
	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).
		Where("is_end = 0").
		Count(&panelGroup.UpcomingTotalCount).Error
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 查询我的待办
	w := service.WorkOrder{
		Classify: 1,
		GinObj:   c,
	}
	result, err = w.PureWorkOrderList()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	panelGroup.MyUpcomingCount = result.(*pagination.Paginator).TotalCount

	app.OK(c, panelGroup, "")
}
