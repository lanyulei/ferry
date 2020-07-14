package service

import (
	"encoding/json"
	"ferry/models/user"
	"ferry/models/workOrder"
	"ferry/pkg/connection"
	"ferry/pkg/pagination"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func WorkOrderList(c *gin.Context, classify int) (result interface{}, err error) {
	type workOrderInfo struct {
		workOrder.Info
		Principals   string `json:"principals"`
		DataClassify int    `json:"data_classify"`
	}
	var (
		workOrderInfoList []workOrderInfo
		principals        string
		userInfo          user.Info
		StateList         []map[string]interface{}
	)

	title := c.DefaultQuery("title", "")
	db := connection.DB.Self.Model(&workOrder.Info{}).Where("title like ?", fmt.Sprintf("%%%v%%", title))

	err = connection.DB.Self.Model(&user.Info{}).Where("id = ?", c.GetInt("userId")).Find(&userInfo).Error
	if err != nil {
		return
	}

	// 获取当前用户信息
	switch classify {
	case 1:
		// 待办工单
		// 1. 个人
		personSelect := fmt.Sprintf("(JSON_CONTAINS(state, JSON_OBJECT('processor', %v)) and JSON_CONTAINS(state, JSON_OBJECT('process_method', 'person')))", c.GetInt("userId"))

		// 2. 小组
		groupList := make([]int, 0)
		err = connection.DB.Self.Model(&user.UserGroup{}).
			Where("user = ?", c.GetInt("userId")).
			Pluck("`group`", &groupList).Error
		if err != nil {
			return
		}
		groupSqlList := make([]string, 0)
		if len(groupList) > 0 {
			for _, group := range groupList {
				groupSqlList = append(groupSqlList, fmt.Sprintf("JSON_CONTAINS(state, JSON_OBJECT('processor', %v))", group))
			}
		} else {
			groupSqlList = append(groupSqlList, fmt.Sprintf("JSON_CONTAINS(state, JSON_OBJECT('processor', 0))"))
		}

		personGroupSelect := fmt.Sprintf(
			"((%v) and %v)",
			strings.Join(groupSqlList, " or "),
			"JSON_CONTAINS(state, JSON_OBJECT('process_method', 'persongroup'))",
		)

		// 3. 部门
		departmentSelect := fmt.Sprintf("(JSON_CONTAINS(state, JSON_OBJECT('processor', %v)) and JSON_CONTAINS(state, JSON_OBJECT('process_method', 'department')))", userInfo.Dept)

		// 4. 变量
		variableSelect := fmt.Sprintf("((%v) or (%v))",
			fmt.Sprintf("JSON_CONTAINS(state, JSON_OBJECT('processor', 1)) and JSON_CONTAINS(state, JSON_OBJECT('process_method', 'variable')) and creator = %v", c.GetInt("userId")),
			fmt.Sprintf("JSON_CONTAINS(state, JSON_OBJECT('processor', 2)) and JSON_CONTAINS(state, JSON_OBJECT('process_method', 'variable')) and creator = %v", userInfo.Dept),
		)

		db = db.Where(fmt.Sprintf("(%v or %v or %v or %v) and is_end = 0", personSelect, personGroupSelect, departmentSelect, variableSelect))
	case 2:
		// 我创建的
		db = db.Where("creator = ?", c.GetInt("userId"))
	case 3:
		// 我相关的
		db = db.Where(fmt.Sprintf("JSON_CONTAINS(related_person, '%v')", c.GetInt("userId")))
	case 4:
	// 所有工单
	default:
		return nil, fmt.Errorf("请确认查询的数据类型是否正确")
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &workOrderInfoList)
	if err != nil {
		err = fmt.Errorf("查询工单列表失败，%v", err.Error())
		return
	}

	for i, w := range *result.(*pagination.Paginator).Data.(*[]workOrderInfo) {
		err = json.Unmarshal(w.State, &StateList)
		if err != nil {
			err = fmt.Errorf("json反序列化失败，%v", err.Error())
			return
		}
		if len(StateList) != 0 {
			processorList := make([]int, 0)
			for _, v := range StateList[0]["processor"].([]interface{}) {
				processorList = append(processorList, int(v.(float64)))
			}
			principals, err = GetPrincipal(processorList, StateList[0]["process_method"].(string))
			if err != nil {
				err = fmt.Errorf("查询处理人名称失败，%v", err.Error())
				return
			}
		}
		workOrderDetails := *result.(*pagination.Paginator).Data.(*[]workOrderInfo)
		workOrderDetails[i].Principals = principals
		workOrderDetails[i].DataClassify = classify
	}

	return result, nil
}
