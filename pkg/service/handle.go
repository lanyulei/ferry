package service

import (
	"encoding/json"
	"errors"
	"ferry/global/orm"
	"ferry/models/base"
	"ferry/models/process"
	"ferry/models/system"
	"ferry/pkg/notify"
	"ferry/tools"
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc : 处理工单
*/

/*
    -- 节点 --
	start: 开始节点
	userTask: 审批节点
	receiveTask: 处理节点
	scriptTask: 任务节点
	end: 结束节点

    -- 网关 --
    exclusiveGateway: 排他网关
    parallelGateway: 并行网关
    inclusiveGateway: 包容网关

*/

type Handle struct {
	cirHistoryList   []process.CirculationHistory
	workOrderId      int
	updateValue      map[string]interface{}
	stateValue       map[string]interface{}
	targetStateValue map[string]interface{}
	WorkOrderData    [][]byte
	workOrderDetails process.WorkOrderInfo
	endHistory       bool
	flowProperties   int
	circulationValue string
	processState     ProcessState
	tx               *gorm.DB
}

// 会签
func (h *Handle) Countersign(c *gin.Context) (err error) {
	var (
		stateList         []map[string]interface{}
		stateIdMap        map[string]interface{}
		currentState      map[string]interface{}
		cirHistoryCount   int
		userInfoList      []system.SysUser
		circulationStatus bool
	)

	err = json.Unmarshal(h.workOrderDetails.State, &stateList)
	if err != nil {
		return
	}

	stateIdMap = make(map[string]interface{})
	for _, v := range stateList {
		stateIdMap[v["id"].(string)] = v["label"]
		if v["id"].(string) == h.stateValue["id"].(string) {
			currentState = v
		}
	}
	userStatusCount := 0
	circulationStatus = false
	for _, cirHistoryValue := range h.cirHistoryList {
		if len(currentState["processor"].([]interface{})) > 1 {
			if _, ok := stateIdMap[cirHistoryValue.Source]; !ok {
				break
			}
		}

		if currentState["process_method"].(string) == "person" {
			// 用户会签
			for _, processor := range currentState["processor"].([]interface{}) {
				if cirHistoryValue.ProcessorId != tools.GetUserId(c) &&
					cirHistoryValue.Source == currentState["id"].(string) &&
					cirHistoryValue.ProcessorId == int(processor.(float64)) {
					cirHistoryCount += 1
				}
			}
			if cirHistoryCount == len(currentState["processor"].([]interface{}))-1 {
				circulationStatus = true
				break
			}
		} else if currentState["process_method"].(string) == "role" || currentState["process_method"].(string) == "department" {
			// 全员处理
			var tmpUserList []system.SysUser
			if h.stateValue["fullHandle"].(bool) {
				db := orm.Eloquent.Model(&system.SysUser{})
				if currentState["process_method"].(string) == "role" {
					db = db.Where("role_id in (?)", currentState["processor"].([]interface{}))
				} else if currentState["process_method"].(string) == "department" {
					db = db.Where("dept_id in (?)", currentState["processor"].([]interface{}))
				}
				err = db.Find(&userInfoList).Error
				if err != nil {
					return
				}
				temp := map[string]struct{}{}
				for _, user := range userInfoList {
					if _, ok := temp[user.Username]; !ok {
						temp[user.Username] = struct{}{}
						tmpUserList = append(tmpUserList, user)
					}
				}
				for _, user := range tmpUserList {
					if cirHistoryValue.Source == currentState["id"].(string) &&
						cirHistoryValue.ProcessorId != tools.GetUserId(c) &&
						cirHistoryValue.ProcessorId == user.UserId {
						userStatusCount += 1
						break
					}
				}
			} else {
				// 普通会签
				for _, processor := range currentState["processor"].([]interface{}) {
					db := orm.Eloquent.Model(&system.SysUser{})
					if currentState["process_method"].(string) == "role" {
						db = db.Where("role_id = ?", processor)
					} else if currentState["process_method"].(string) == "department" {
						db = db.Where("dept_id = ?", processor)
					}
					err = db.Find(&userInfoList).Error
					if err != nil {
						return
					}
					for _, user := range userInfoList {
						if user.UserId != tools.GetUserId(c) &&
							cirHistoryValue.Source == currentState["id"].(string) &&
							cirHistoryValue.ProcessorId == user.UserId {
							userStatusCount += 1
							break
						}
					}
				}
			}
			if h.stateValue["fullHandle"].(bool) {
				if userStatusCount == len(tmpUserList)-1 {
					circulationStatus = true
				}
			} else {
				if userStatusCount == len(currentState["processor"].([]interface{}))-1 {
					circulationStatus = true
				}
			}
		}
	}
	if circulationStatus {
		h.endHistory = true
		err = h.circulation()
		if err != nil {
			return
		}
	}
	return
}

// 工单跳转
func (h *Handle) circulation() (err error) {
	var (
		stateValue []byte
	)

	stateList := make([]interface{}, 0)
	for _, v := range h.updateValue["state"].([]map[string]interface{}) {
		stateList = append(stateList, v)
	}
	err = GetVariableValue(stateList, h.workOrderDetails.Creator)
	if err != nil {
		return
	}

	stateValue, err = json.Marshal(h.updateValue["state"])
	if err != nil {
		return
	}

	err = h.tx.Model(&process.WorkOrderInfo{}).
		Where("id = ?", h.workOrderId).
		Updates(map[string]interface{}{
			"state":          stateValue,
			"related_person": h.updateValue["related_person"],
		}).Error
	if err != nil {
		h.tx.Rollback()
		return
	}

	// 如果是跳转到结束节点，则需要修改节点状态
	if h.targetStateValue["clazz"] == "end" {
		err = h.tx.Model(&process.WorkOrderInfo{}).
			Where("id = ?", h.workOrderId).
			Update("is_end", 1).Error
		if err != nil {
			h.tx.Rollback()
			return
		}
	}

	return
}

// 条件判断
func (h *Handle) ConditionalJudgment(condExpr map[string]interface{}) (result bool, err error) {
	var (
		condExprOk    bool
		condExprValue interface{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case string:
				err = errors.New(e)
			case error:
				err = e
			default:
				err = errors.New("未知错误")
			}
			return
		}
	}()

	for _, data := range h.WorkOrderData {
		var formData map[string]interface{}
		err = json.Unmarshal(data, &formData)
		if err != nil {
			return
		}
		if condExprValue, condExprOk = formData[condExpr["key"].(string)]; condExprOk {
			break
		}
	}

	if condExprValue == nil {
		err = errors.New("未查询到对应的表单数据。")
		return
	}

	// todo 待优化
	switch reflect.TypeOf(condExprValue).String() {
	case "string":
		switch condExpr["sign"] {
		case "==":
			if condExprValue.(string) == condExpr["value"].(string) {
				result = true
			}
		case "!=":
			if condExprValue.(string) != condExpr["value"].(string) {
				result = true
			}
		case ">":
			if condExprValue.(string) > condExpr["value"].(string) {
				result = true
			}
		case ">=":
			if condExprValue.(string) >= condExpr["value"].(string) {
				result = true
			}
		case "<":
			if condExprValue.(string) < condExpr["value"].(string) {
				result = true
			}
		case "<=":
			if condExprValue.(string) <= condExpr["value"].(string) {
				result = true
			}
		default:
			err = errors.New("目前仅支持6种常规判断类型，包括（等于、不等于、大于、大于等于、小于、小于等于）")
		}
	case "float64":
		switch condExpr["sign"] {
		case "==":
			if condExprValue.(float64) == condExpr["value"].(float64) {
				result = true
			}
		case "!=":
			if condExprValue.(float64) != condExpr["value"].(float64) {
				result = true
			}
		case ">":
			if condExprValue.(float64) > condExpr["value"].(float64) {
				result = true
			}
		case ">=":
			if condExprValue.(float64) >= condExpr["value"].(float64) {
				result = true
			}
		case "<":
			if condExprValue.(float64) < condExpr["value"].(float64) {
				result = true
			}
		case "<=":
			if condExprValue.(float64) <= condExpr["value"].(float64) {
				result = true
			}
		default:
			err = errors.New("目前仅支持6种常规判断类型，包括（等于、不等于、大于、大于等于、小于、小于等于）")
		}
	default:
		err = errors.New("条件判断目前仅支持字符串、整型。")
	}

	return
}

// 并行网关，确认其他节点是否完成
func (h *Handle) completeAllParallel(target string) (statusOk bool, err error) {
	var (
		stateList []map[string]interface{}
	)

	err = json.Unmarshal(h.workOrderDetails.State, &stateList)
	if err != nil {
		err = fmt.Errorf("反序列化失败，%v", err.Error())
		return
	}

continueHistoryTag:
	for _, v := range h.cirHistoryList {
		status := false
		for i, s := range stateList {
			if v.Source == s["id"].(string) && v.Target == target {
				status = true
				stateList = append(stateList[:i], stateList[i+1:]...)
				continue continueHistoryTag
			}
		}
		if !status {
			break
		}
	}

	if len(stateList) == 1 && stateList[0]["id"].(string) == h.stateValue["id"] {
		statusOk = true
	}

	return
}

func (h *Handle) commonProcessing(c *gin.Context) (err error) {
	// 如果是拒绝的流转则直接跳转
	if h.flowProperties == 0 {
		err = h.circulation()
		if err != nil {
			err = fmt.Errorf("工单跳转失败，%v", err.Error())
		}
		return
	}

	// 会签
	if h.stateValue["assignValue"] != nil && len(h.stateValue["assignValue"].([]interface{})) > 0 {
		if isCounterSign, ok := h.stateValue["isCounterSign"]; ok {
			if isCounterSign.(bool) {
				h.endHistory = false
				err = h.Countersign(c)
				if err != nil {
					return
				}
			} else {
				err = h.circulation()
				if err != nil {
					return
				}
			}
		} else {
			err = h.circulation()
			if err != nil {
				return
			}
		}
	} else {
		err = h.circulation()
		if err != nil {
			return
		}
	}
	return
}

func (h *Handle) HandleWorkOrder(
	c *gin.Context,
	workOrderId int,
	tasks []string,
	targetState string,
	sourceState string,
	circulationValue string,
	flowProperties int,
	remarks string,
	tpls []map[string]interface{},
	isExecTask bool,
) (err error) {
	h.workOrderId = workOrderId
	h.flowProperties = flowProperties
	h.endHistory = true

	var (
		execTasks          []string
		relatedPersonList  []int
		cirHistoryValue    []process.CirculationHistory
		cirHistoryData     process.CirculationHistory
		costDurationValue  int64
		sourceEdges        []map[string]interface{}
		targetEdges        []map[string]interface{}
		condExprStatus     bool
		relatedPersonValue []byte
		parallelStatusOk   bool
		processInfo        process.Info
		currentUserInfo    system.SysUser
		applyUserInfo      system.SysUser
		sendToUserList     []system.SysUser
		noticeList         []int
		sendSubject        string = "您有一条待办工单，请及时处理"
		sendDescription    string = "您有一条待办工单请及时处理，工单描述如下"
		paramsValue        struct {
			Id       int           `json:"id"`
			Title    string        `json:"title"`
			Priority int           `json:"priority"`
			FormData []interface{} `json:"form_data"`
		}
	)

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case string:
				err = errors.New(e)
			case error:
				err = e
			default:
				err = errors.New("未知错误")
			}
			return
		}
	}()

	// 获取工单信息
	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).Where("id = ?", workOrderId).Find(&h.workOrderDetails).Error
	if err != nil {
		return
	}

	// 查询工单创建人信息
	err = orm.Eloquent.Model(&system.SysUser{}).Where("user_id = ?", h.workOrderDetails.Creator).Find(&applyUserInfo).Error
	if err != nil {
		return
	}

	// 获取流程信息
	err = orm.Eloquent.Model(&process.Info{}).Where("id = ?", h.workOrderDetails.Process).Find(&processInfo).Error
	if err != nil {
		return
	}
	err = json.Unmarshal(processInfo.Structure, &h.processState.Structure)
	if err != nil {
		return
	}

	// 获取当前节点
	h.stateValue, err = h.processState.GetNode(sourceState)
	if err != nil {
		return
	}

	// 目标状态
	h.targetStateValue, err = h.processState.GetNode(targetState)
	if err != nil {
		return
	}

	// 获取工单数据
	err = orm.Eloquent.Model(&process.TplData{}).
		Where("work_order = ?", workOrderId).
		Pluck("form_data", &h.WorkOrderData).Error
	if err != nil {
		return
	}

	// 根据处理人查询出需要会签的条数
	err = orm.Eloquent.Model(&process.CirculationHistory{}).
		Where("work_order = ?", workOrderId).
		Order("id desc").
		Find(&h.cirHistoryList).Error
	if err != nil {
		return
	}

	err = json.Unmarshal(h.workOrderDetails.RelatedPerson, &relatedPersonList)
	if err != nil {
		return
	}
	relatedPersonStatus := false
	for _, r := range relatedPersonList {
		if r == tools.GetUserId(c) {
			relatedPersonStatus = true
			break
		}
	}
	if !relatedPersonStatus {
		relatedPersonList = append(relatedPersonList, tools.GetUserId(c))
	}

	relatedPersonValue, err = json.Marshal(relatedPersonList)
	if err != nil {
		return
	}

	h.updateValue = map[string]interface{}{
		"related_person": relatedPersonValue,
	}

	// 开启事务
	h.tx = orm.Eloquent.Begin()

	stateValue := map[string]interface{}{
		"label": h.targetStateValue["label"].(string),
		"id":    h.targetStateValue["id"].(string),
	}

	sourceEdges, err = h.processState.GetEdge(h.targetStateValue["id"].(string), "source")
	if err != nil {
		return
	}

	switch h.targetStateValue["clazz"] {
	case "exclusiveGateway": // 排他网关
	breakTag:
		for _, edge := range sourceEdges {
			edgeCondExpr := make([]map[string]interface{}, 0)
			err = json.Unmarshal([]byte(edge["conditionExpression"].(string)), &edgeCondExpr)
			if err != nil {
				return
			}
			for _, condExpr := range edgeCondExpr {
				// 条件判断
				condExprStatus, err = h.ConditionalJudgment(condExpr)
				if err != nil {
					return
				}
				if condExprStatus {
					// 进行节点跳转
					h.targetStateValue, err = h.processState.GetNode(edge["target"].(string))
					if err != nil {
						return
					}

					if h.targetStateValue["clazz"] == "userTask" || h.targetStateValue["clazz"] == "receiveTask" {
						if h.targetStateValue["assignValue"] == nil || h.targetStateValue["assignType"] == "" {
							err = errors.New("处理人不能为空")
							return
						}
					}

					h.updateValue["state"] = []map[string]interface{}{{
						"id":             h.targetStateValue["id"].(string),
						"label":          h.targetStateValue["label"],
						"processor":      h.targetStateValue["assignValue"],
						"process_method": h.targetStateValue["assignType"],
					}}
					err = h.commonProcessing(c)
					if err != nil {
						err = fmt.Errorf("流程流程跳转失败，%v", err.Error())
						return
					}

					break breakTag
				}
			}
		}
		if !condExprStatus {
			err = errors.New("所有流转均不符合条件，请确认。")
			return
		}
	case "parallelGateway": // 并行/聚合网关
		// 入口，判断
		targetEdges, err = h.processState.GetEdge(h.targetStateValue["id"].(string), "target")
		if err != nil {
			err = fmt.Errorf("查询流转信息失败，%v", err.Error())
			return
		}

		if len(sourceEdges) > 0 {
			h.targetStateValue, err = h.processState.GetNode(sourceEdges[0]["target"].(string))
			if err != nil {
				return
			}
		} else {
			err = errors.New("并行网关流程不正确")
			return
		}

		if len(sourceEdges) > 1 && len(targetEdges) == 1 {
			// 入口
			h.updateValue["state"] = make([]map[string]interface{}, 0)
			for _, edge := range sourceEdges {
				targetStateValue, err := h.processState.GetNode(edge["target"].(string))
				if err != nil {
					return err
				}
				h.updateValue["state"] = append(h.updateValue["state"].([]map[string]interface{}), map[string]interface{}{
					"id":             edge["target"].(string),
					"label":          targetStateValue["label"],
					"processor":      targetStateValue["assignValue"],
					"process_method": targetStateValue["assignType"],
				})
			}
			err = h.circulation()
			if err != nil {
				err = fmt.Errorf("工单跳转失败，%v", err.Error())
				return
			}
		} else if len(sourceEdges) == 1 && len(targetEdges) > 1 {
			// 出口
			parallelStatusOk, err = h.completeAllParallel(sourceEdges[0]["target"].(string))
			if err != nil {
				err = fmt.Errorf("并行检测失败，%v", err.Error())
				return
			}
			if parallelStatusOk {
				h.endHistory = true
				endAssignValue, ok := h.targetStateValue["assignValue"]
				if !ok {
					endAssignValue = []int{}
				}

				endAssignType, ok := h.targetStateValue["assignType"]
				if !ok {
					endAssignType = ""
				}

				h.updateValue["state"] = []map[string]interface{}{{
					"id":             h.targetStateValue["id"].(string),
					"label":          h.targetStateValue["label"],
					"processor":      endAssignValue,
					"process_method": endAssignType,
				}}
				err = h.circulation()
				if err != nil {
					err = fmt.Errorf("工单跳转失败，%v", err.Error())
					return
				}
			} else {
				h.endHistory = false
			}

		} else {
			err = errors.New("并行网关流程不正确")
			return
		}
	// 包容网关
	case "inclusiveGateway":
		return
	case "start":
		stateValue["processor"] = []int{h.workOrderDetails.Creator}
		stateValue["process_method"] = "person"
		h.updateValue["state"] = []map[string]interface{}{stateValue}
		err = h.circulation()
		if err != nil {
			return
		}
	case "userTask":
		stateValue["processor"] = h.targetStateValue["assignValue"].([]interface{})
		stateValue["process_method"] = h.targetStateValue["assignType"].(string)
		h.updateValue["state"] = []map[string]interface{}{stateValue}
		err = h.commonProcessing(c)
		if err != nil {
			return
		}
	case "receiveTask":
		stateValue["processor"] = h.targetStateValue["assignValue"].([]interface{})
		stateValue["process_method"] = h.targetStateValue["assignType"].(string)
		h.updateValue["state"] = []map[string]interface{}{stateValue}
		err = h.commonProcessing(c)
		if err != nil {
			return
		}
	case "scriptTask":
		stateValue["processor"] = []int{}
		stateValue["process_method"] = ""
		h.updateValue["state"] = []map[string]interface{}{stateValue}
	case "end":
		stateValue["processor"] = []int{}
		stateValue["process_method"] = ""
		h.updateValue["state"] = []map[string]interface{}{stateValue}
		err = h.commonProcessing(c)
		if err != nil {
			return
		}
	}

	// 更新表单数据
	for _, t := range tpls {
		var (
			tplValue []byte
		)
		tplValue, err = json.Marshal(t["tplValue"])
		if err != nil {
			h.tx.Rollback()
			return
		}

		paramsValue.FormData = append(paramsValue.FormData, t["tplValue"])

		// 是否可写，只有可写的模版可以更新数据
		updateStatus := false
		if h.stateValue["clazz"].(string) == "start" {
			updateStatus = true
		} else if writeTplList, writeOK := h.stateValue["writeTpls"]; writeOK {
		tplListTag:
			for _, writeTplId := range writeTplList.([]interface{}) {
				if writeTplId == t["tplId"] { // 可写
					// 是否隐藏，隐藏的模版无法修改数据
					if hideTplList, hideOK := h.stateValue["hideTpls"]; hideOK {
						if hideTplList != nil && len(hideTplList.([]interface{})) > 0 {
							for _, hideTplId := range hideTplList.([]interface{}) {
								if hideTplId == t["tplId"] { // 隐藏的
									updateStatus = false
									break tplListTag
								} else {
									updateStatus = true
								}
							}
						} else {
							updateStatus = true
						}
					} else {
						updateStatus = true
					}
				}
			}
		} else {
			// 不可写
			updateStatus = false
		}
		if updateStatus {
			err = h.tx.Model(&process.TplData{}).Where("id = ?", t["tplDataId"]).Update("form_data", tplValue).Error
			if err != nil {
				h.tx.Rollback()
				return
			}
		}
	}

	// 流转历史写入
	err = orm.Eloquent.Model(&cirHistoryValue).
		Where("work_order = ?", workOrderId).
		Find(&cirHistoryValue).
		Order("create_time desc").Error
	if err != nil {
		h.tx.Rollback()
		return
	}
	for _, t := range cirHistoryValue {
		if t.Source != h.stateValue["id"] {
			costDuration := time.Since(t.CreatedAt.Time)
			costDurationValue = int64(costDuration) / 1000 / 1000 / 1000
		}
	}

	// 获取当前用户信息
	err = orm.Eloquent.Model(&currentUserInfo).
		Where("user_id = ?", tools.GetUserId(c)).
		Find(&currentUserInfo).Error
	if err != nil {
		return
	}

	cirHistoryData = process.CirculationHistory{
		Model:        base.Model{},
		Title:        h.workOrderDetails.Title,
		WorkOrder:    h.workOrderDetails.Id,
		State:        h.stateValue["label"].(string),
		Source:       h.stateValue["id"].(string),
		Target:       h.targetStateValue["id"].(string),
		Circulation:  circulationValue,
		Processor:    currentUserInfo.NickName,
		ProcessorId:  tools.GetUserId(c),
		Status:       flowProperties,
		CostDuration: costDurationValue,
		Remarks:      remarks,
	}
	err = h.tx.Create(&cirHistoryData).Error
	if err != nil {
		h.tx.Rollback()
		return
	}

	// 获取流程通知类型列表
	err = json.Unmarshal(processInfo.Notice, &noticeList)
	if err != nil {
		return
	}

	// 获取需要抄送的邮件
	emailCCList := make([]string, 0)
	if h.stateValue["cc"] != nil && len(h.stateValue["cc"].([]interface{})) > 0 {
		err = orm.Eloquent.Model(&system.SysUser{}).
			Where("user_id in (?)", h.stateValue["cc"]).
			Pluck("email", &emailCCList).Error
		if err != nil {
			err = errors.New("查询邮件抄送人失败")
			return
		}
	}

	bodyData := notify.BodyData{
		SendTo: map[string]interface{}{
			"userList": sendToUserList,
		},
		EmailCcTo:   emailCCList,
		Subject:     sendSubject,
		Description: sendDescription,
		Classify:    noticeList,
		ProcessId:   h.workOrderDetails.Process,
		Id:          h.workOrderDetails.Id,
		Title:       h.workOrderDetails.Title,
		Creator:     applyUserInfo.NickName,
		Priority:    h.workOrderDetails.Priority,
		CreatedAt:   h.workOrderDetails.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 判断目标是否是结束节点
	if h.targetStateValue["clazz"] == "end" && h.endHistory == true {
		sendSubject = "您的工单已处理完成"
		sendDescription = "您的工单已处理完成，工单描述如下"
		err = h.tx.Create(&process.CirculationHistory{
			Model:       base.Model{},
			Title:       h.workOrderDetails.Title,
			WorkOrder:   h.workOrderDetails.Id,
			State:       h.targetStateValue["label"].(string),
			Source:      h.targetStateValue["id"].(string),
			Processor:   currentUserInfo.NickName,
			ProcessorId: tools.GetUserId(c),
			Circulation: "结束",
			Remarks:     "工单已结束",
			Status:      2, // 其他状态
		}).Error
		if err != nil {
			h.tx.Rollback()
			return
		}
		if len(noticeList) > 0 {
			// 查询工单创建人信息
			err = h.tx.Model(&system.SysUser{}).
				Where("user_id = ?", h.workOrderDetails.Creator).
				Find(&sendToUserList).Error
			if err != nil {
				return
			}

			bodyData.SendTo = map[string]interface{}{
				"userList": sendToUserList,
			}
			bodyData.Subject = sendSubject
			bodyData.Description = sendDescription

			// 发送通知
			go func(bodyData notify.BodyData) {
				err = bodyData.SendNotify()
				if err != nil {
					return
				}
			}(bodyData)
		}
	}

	h.tx.Commit() // 提交事务

	// 发送通知
	if len(noticeList) > 0 {
		stateList := make([]interface{}, 0)
		for _, v := range h.updateValue["state"].([]map[string]interface{}) {
			stateList = append(stateList, v)
		}
		sendToUserList, err = GetPrincipalUserInfo(stateList, h.workOrderDetails.Creator)
		if err != nil {
			return
		}

		bodyData.SendTo = map[string]interface{}{
			"userList": sendToUserList,
		}
		bodyData.Subject = sendSubject
		bodyData.Description = sendDescription

		// 发送通知
		go func(bodyData notify.BodyData) {
			err = bodyData.SendNotify()
			if err != nil {
				return
			}
		}(bodyData)
	}

	if isExecTask {
		// 执行流程公共任务及节点任务
		if h.stateValue["task"] != nil {
			for _, task := range h.stateValue["task"].([]interface{}) {
				tasks = append(tasks, task.(string))
			}
		}
	continueTag:
		for _, task := range tasks {
			for _, t := range execTasks {
				if t == task {
					continue continueTag
				}
			}
			execTasks = append(execTasks, task)
		}

		paramsValue.Id = h.workOrderDetails.Id
		paramsValue.Title = h.workOrderDetails.Title
		paramsValue.Priority = h.workOrderDetails.Priority
		params, err := json.Marshal(paramsValue)
		if err != nil {
			return err
		}
		go ExecTask(execTasks, string(params))
	}
	return
}
