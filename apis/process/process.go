package process

import (
	"errors"
	"ferry/global/orm"
	"ferry/models/process"
	process2 "ferry/models/process"
	"ferry/pkg/pagination"
	"ferry/tools"
	"ferry/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
 @Author : lanyulei
*/

// 流程列表
func ProcessList(c *gin.Context) {
	var (
		err         error
		processList []*struct {
			process2.Info
			CreateUser   string `json:"create_user"`
			CreateName   string `json:"create_name"`
			ClassifyName string `json:"classify_name"`
		}
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&process2.Info{}).
		Joins("left join sys_user on sys_user.user_id = p_process_info.creator").
		Joins("left join p_process_classify on p_process_classify.id = p_process_info.classify").
		Select("p_process_info.id, p_process_info.create_time, p_process_info.update_time, p_process_info.name, p_process_info.creator, p_process_classify.name as classify_name, sys_user.username as create_user, sys_user.nick_name as create_name")

	result, err := pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &processList, SearchParams, "p_process_info")

	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询流程模版失败，%v", err.Error()))
		return
	}
	app.OK(c, result, "查询流程列表成功")
}

// 创建流程
func CreateProcess(c *gin.Context) {
	var (
		err          error
		processValue process2.Info
		processCount int
	)

	err = c.ShouldBind(&processValue)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("参数绑定失败，%v", err.Error()))
		return
	}

	// 确定流程名称是否重复
	err = orm.Eloquent.Model(&processValue).
		Where("name = ?", processValue.Name).
		Count(&processCount).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("流程信息查询失败，%v", err.Error()))
		return
	}
	if processCount > 0 {
		app.Error(c, -1, errors.New("流程名称重复"), "")
		return
	}

	processValue.Creator = tools.GetUserId(c)

	err = orm.Eloquent.Create(&processValue).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("创建流程失败，%v", err.Error()))
		return
	}

	app.OK(c, processValue, "流程创建成功")
}

// 更新流程
func UpdateProcess(c *gin.Context) {
	var (
		err          error
		processValue process2.Info
	)

	err = c.ShouldBind(&processValue)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	err = orm.Eloquent.Model(&process2.Info{}).
		Where("id = ?", processValue.Id).
		Updates(map[string]interface{}{
			"name":      processValue.Name,
			"structure": processValue.Structure,
			"tpls":      processValue.Tpls,
			"classify":  processValue.Classify,
			"task":      processValue.Task,
			"notice":    processValue.Notice,
			"icon":      processValue.Icon,
			"remarks":   processValue.Remarks,
		}).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("更新流程信息失败，%v", err.Error()))
		return
	}

	app.OK(c, processValue, "更新流程信息成功")
}

// 删除流程
func DeleteProcess(c *gin.Context) {
	processId := c.DefaultQuery("processId", "")
	if processId == "" {
		app.Error(c, -1, errors.New("参数不正确，请确定参数processId是否传递"), "")
		return
	}

	err := orm.Eloquent.Delete(process2.Info{}, "id = ?", processId).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("删除流程失败, %v", err.Error()))
		return
	}
	app.OK(c, "", "删除流程成功")
}

// 流程详情
func ProcessDetails(c *gin.Context) {
	processId := c.DefaultQuery("processId", "")
	if processId == "" {
		app.Error(c, -1, errors.New("参数不正确，请确定参数processId是否传递"), "")
		return
	}

	var processValue process2.Info
	err := orm.Eloquent.Model(&processValue).
		Where("id = ?", processId).
		Find(&processValue).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询流程详情失败, %v", err.Error()))
		return
	}

	app.OK(c, processValue, "查询流程详情成功")
}

// 分类流程列表
func ClassifyProcessList(c *gin.Context) {
	var (
		err            error
		classifyIdList []int
		classifyList   []*struct {
			process2.Classify
			ProcessList []*process2.Info `json:"process_list"`
		}
	)

	processName := c.DefaultQuery("name", "")
	if processName == "" {
		err = orm.Eloquent.Model(&process2.Classify{}).Find(&classifyList).Error
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("获取分类列表失败，%v", err.Error()))
			return
		}
	} else {
		err = orm.Eloquent.Model(&process2.Info{}).
			Where("name LIKE ?", fmt.Sprintf("%%%v%%", processName)).
			Pluck("distinct classify", &classifyIdList).Error
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
			return
		}

		err = orm.Eloquent.Model(&process2.Classify{}).
			Where("id in (?)", classifyIdList).
			Find(&classifyList).Error
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("获取分类失败，%v", err.Error()))
			return
		}
	}

	for _, item := range classifyList {
		err = orm.Eloquent.Model(&process2.Info{}).
			Where("classify = ? and name LIKE ?", item.Id, fmt.Sprintf("%%%v%%", processName)).
			Select("id, create_time, update_time, name, icon, remarks").
			Find(&item.ProcessList).Error
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("获取流程失败，%v", err.Error()))
			return
		}
	}

	app.OK(c, classifyList, "成功获取数据")
}

// 克隆流程
func CloneProcess(c *gin.Context) {
	var (
		err  error
		id   string
		info process.Info
	)

	id = c.Param("id")

	err = orm.Eloquent.Find(&info, id).Error
	if err != nil {
		app.Error(c, -1, err, "查询流程数据失败")
		return
	}

	err = orm.Eloquent.Create(&process.Info{
		Name:        info.Name + "-copy",
		Icon:        info.Icon,
		Structure:   info.Structure,
		Classify:    info.Classify,
		Tpls:        info.Tpls,
		Task:        info.Task,
		SubmitCount: info.SubmitCount,
		Creator:     tools.GetUserId(c),
		Notice:      info.Notice,
		Remarks:     info.Remarks,
	}).Error
	if err != nil {
		app.Error(c, -1, err, "克隆流程失败")
		return
	}

	app.OK(c, nil, "")
}
