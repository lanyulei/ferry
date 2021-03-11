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
		err   error
		count map[string]int // 工单数量统计
		ranks []service.Ranks
	)

	statistics := service.Statistics{
		StartTime: "",
		EndTime:   "",
	}

	// 查询工单类型数据统计
	count, err = statistics.WorkOrderCount(c)
	if err != nil {
		app.Error(c, -1, err, "查询工单类型数据统计失败")
		return
	}

	// 查询工单数据排名
	ranks, err = statistics.WorkOrderRanks()
	if err != nil {
		app.Error(c, -1, err, "查询提交工单排名数据失败")
		return
	}

	app.OK(c, map[string]interface{}{
		"count": count,
		"ranks": ranks,
	}, "")
}
