package process

//import (
//	process2 "ferry/models/process"
//	"ferry/pkg/connection"
//	"ferry/pkg/pagination"
//	"ferry/pkg/response/code"
//	"fmt"
//
//	"github.com/gin-gonic/gin"
//)
//
///*
//  @Author : lanyulei
//*/
//
//// 流程列表
//func ProcessList(c *gin.Context) {
//	type processValue struct {
//		process2.Info
//		CreateUser   string `json:"create_user"`
//		CreateName   string `json:"create_name"`
//		ClassifyName string `json:"classify_name"`
//	}
//
//	var (
//		err         error
//		processList []*processValue
//	)
//
//	SearchParams := map[string]map[string]interface{}{
//		"like": pagination.RequestParams(c),
//	}
//
//	db := connection.DB.Self.
//		Model(&process2.Info{}).
//		Joins("left join user_info on user_info.id = p_process_info.creator").
//		Joins("left join p_process_classify on p_process_classify.id = p_process_info.classify").
//		Select("p_process_info.id, p_process_info.create_time, p_process_info.update_time, p_process_info.name, p_process_info.creator, p_process_classify.name as classify_name, user_info.username as create_user, user_info.nickname as create_name").
//		Where("p_process_info.`delete_time` IS NULL")
//
//	result, err := pagination.Paging(&pagination.Param{
//		C:  c,
//		DB: db,
//	}, &processList, SearchParams, "p_process_info")
//
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询流程列表失败，%v", err.Error()))
//		return
//	}
//	Response(c, nil, result, "")
//}
//
//// 创建流程
//func CreateProcess(c *gin.Context) {
//	var (
//		err          error
//		processValue process2.Info
//		processCount int
//	)
//
//	err = c.ShouldBind(&processValue)
//	if err != nil {
//		Response(c, code.BindError, nil, err.Error())
//		return
//	}
//
//	// 确定修改的分类是否存在
//	err = connection.DB.Self.Model(&processValue).
//		Where("name = ?", processValue.Name).
//		Count(&processCount).Error
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询流程数量失败，%v", err.Error()))
//		return
//	}
//	if processCount > 0 {
//		Response(c, code.InternalServerError, nil, "流程名称出现重复，请换一个名称")
//		return
//	}
//
//	processValue.Creator = c.GetInt("userId")
//
//	err = connection.DB.Self.Create(&processValue).Error
//	if err != nil {
//		Response(c, code.CreateError, nil, fmt.Sprintf("创建流程失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 更新流程
//func UpdateProcess(c *gin.Context) {
//	var (
//		err          error
//		processValue process2.Info
//	)
//
//	err = c.ShouldBind(&processValue)
//	if err != nil {
//		Response(c, code.BindError, nil, err.Error())
//		return
//	}
//
//	err = connection.DB.Self.Model(&process2.Info{}).
//		Where("id = ?", processValue.Id).
//		Updates(map[string]interface{}{
//			"name":      processValue.Name,
//			"structure": processValue.Structure,
//			"tpls":      processValue.Tpls,
//			"classify":  processValue.Classify,
//			"task":      processValue.Task,
//		}).Error
//	if err != nil {
//		Response(c, code.UpdateError, nil, fmt.Sprintf("更新流程信息失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 删除流程
//func DeleteProcess(c *gin.Context) {
//	processId := c.DefaultQuery("processId", "")
//	if processId == "" {
//		Response(c, code.InternalServerError, nil, "参数不正确，请确定参数processId是否传递")
//		return
//	}
//
//	err := connection.DB.Self.Delete(process2.Info{}, "id = ?", processId).Error
//	if err != nil {
//		Response(c, code.DeleteError, nil, fmt.Sprintf("删除流程失败, %v", err.Error()))
//		return
//	}
//	Response(c, nil, nil, "")
//}
//
//// 流程详情
//func ProcessDetails(c *gin.Context) {
//	processId := c.DefaultQuery("processId", "")
//	if processId == "" {
//		Response(c, code.InternalServerError, nil, "参数不正确，请确定参数processId是否传递")
//		return
//	}
//
//	var processValue process2.Info
//	err := connection.DB.Self.Model(&processValue).
//		Where("id = ?", processId).
//		Find(&processValue).Error
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询流程详情失败, %v", err.Error()))
//		return
//	}
//
//	Response(c, nil, processValue, "")
//}
//
//// 分类流程列表
//func ClassifyProcessList(c *gin.Context) {
//	type classifyProcess struct {
//		process2.Classify
//		ProcessList []*process2.Info `json:"process_list"`
//	}
//
//	var (
//		err          error
//		classifyList []*classifyProcess
//	)
//
//	processName := c.DefaultQuery("name", "")
//	if processName == "" {
//		err = connection.DB.Self.Model(&process2.Classify{}).Find(&classifyList).Error
//		if err != nil {
//			Response(c, code.SelectError, nil, fmt.Sprintf("获取分类列表失败，%v", err.Error()))
//			return
//		}
//	} else {
//		var classifyIdList []int
//		err = connection.DB.Self.Model(&process2.Info{}).
//			Where("name LIKE ?", fmt.Sprintf("%%%v%%", processName)).
//			Pluck("distinct classify", &classifyIdList).Error
//		if err != nil {
//			Response(c, code.SelectError, nil, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//
//		err = connection.DB.Self.Model(&process2.Classify{}).
//			Where("id in (?)", classifyIdList).
//			Find(&classifyList).Error
//		if err != nil {
//			Response(c, code.SelectError, nil, fmt.Sprintf("获取分类失败，%v", err.Error()))
//			return
//		}
//	}
//
//	for _, item := range classifyList {
//		err = connection.DB.Self.Model(&process2.Info{}).
//			Where("classify = ?", item.Id).
//			Select("id, create_time, update_time, name").
//			Find(&item.ProcessList).Error
//		if err != nil {
//			Response(c, code.SelectError, nil, fmt.Sprintf("获取流程失败，%v", err.Error()))
//			return
//		}
//	}
//
//	Response(c, nil, classifyList, "")
//}
