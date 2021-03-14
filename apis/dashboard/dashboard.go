package dashboard

import (
	"ferry/pkg/service"
	"ferry/tools/app"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func InitData(c *gin.Context) {
	var (
		err       error
		count     map[string]int // 工单数量统计
		ranks     []service.Ranks
		submit    map[string][]interface{}
		startTime string
		endTime   string
		handle    interface{}
		period    interface{}
	)

	startTime = c.DefaultQuery("start_time", "")
	endTime = c.DefaultQuery("end_time", "")

	if startTime == "" || endTime == "" {
		// 默认为最近7天的数据
		startTime = fmt.Sprintf("%s 00:00:00", time.Now().AddDate(0, 0, -6).Format("2006-01-02"))
		endTime = fmt.Sprintf("%s 23:59:59", time.Now().Format("2006-01-02"))
	} else {
		if strings.Contains(startTime, "T") && strings.Contains(endTime, "T") {
			startTime = fmt.Sprintf("%s 00:00:00", strings.Split(startTime, "T")[0])
			endTime = fmt.Sprintf("%s 23:59:59", strings.Split(endTime, "T")[0])
		}
	}

	statistics := service.Statistics{
		StartTime: startTime,
		EndTime:   endTime,
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

	// 工单提交数据统计
	submit, err = statistics.DateRangeStatistics()
	if err != nil {
		app.Error(c, -1, err, "工单提交数据统计失败")
		return
	}

	// 处理工单人员排行榜
	handle, err = statistics.HandlePersonRank()
	if err != nil {
		app.Error(c, -1, err, "查询处理工单人员排行失败")
		return
	}

	// 工单处理耗时排行榜
	period, err = statistics.HandlePeriodRank()
	if err != nil {
		app.Error(c, -1, err, "查询工单处理耗时排行失败")
		return
	}

	app.OK(c, map[string]interface{}{
		"count":  count,
		"ranks":  ranks,
		"submit": submit,
		"handle": handle,
		"period": period,
	}, "")
}
