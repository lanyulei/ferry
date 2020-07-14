package service

import (
	"encoding/json"
	"ferry/models/process"
	"ferry/models/user"
	"ferry/models/workOrder"
	"ferry/pkg/connection"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func JudgeUserAuthority(c *gin.Context, workOrderId int, currentState string) (status bool, err error) {
	/*
		person 人员
		persongroup 人员组
		department 部门
		variable 变量
	*/
	var (
		workOrderInfo     workOrder.Info
		userInfo          user.Info
		userDept          user.Dept
		cirHistoryList    []workOrder.CirculationHistory
		stateValue        map[string]interface{}
		processInfo       process.Info
		processState      ProcessState
		currentStateList  []map[string]interface{}
		currentStateValue map[string]interface{}
	)
	// 获取工单信息
	err = connection.DB.Self.Model(&workOrderInfo).
		Where("id = ?", workOrderId).
		Find(&workOrderInfo).Error
	if err != nil {
		return
	}

	// 获取流程信息
	err = connection.DB.Self.Model(&process.Info{}).Where("id = ?", workOrderInfo.Process).Find(&processInfo).Error
	if err != nil {
		return
	}
	err = json.Unmarshal(processInfo.Structure, &processState.Structure)
	if err != nil {
		return
	}

	stateValue, err = processState.GetNode(currentState)
	if err != nil {
		return
	}

	err = json.Unmarshal(workOrderInfo.State, &currentStateList)
	if err != nil {
		return
	}

	for _, v := range currentStateList {
		if v["id"].(string) == currentState {
			currentStateValue = v
			break
		}
	}

	// 会签
	if currentStateValue["processor"] != nil && len(currentStateValue["processor"].([]interface{})) > 1 {
		if isCounterSign, ok := stateValue["isCounterSign"]; ok {
			if isCounterSign.(bool) {
				err = connection.DB.Self.Model(&workOrder.CirculationHistory{}).
					Where("work_order = ?", workOrderId).
					Order("id desc").
					Find(&cirHistoryList).Error
				if err != nil {
					return
				}
				for _, cirHistoryValue := range cirHistoryList {
					if cirHistoryValue.Source != stateValue["id"] {
						break
					}
					if cirHistoryValue.Source == stateValue["id"] && cirHistoryValue.ProcessorId == c.GetInt("userId") {
						return
					}
				}
			}
		}
	}

	switch currentStateValue["process_method"].(string) {
	case "person":
		for _, processorValue := range currentStateValue["processor"].([]interface{}) {
			if int(processorValue.(float64)) == c.GetInt("userId") {
				status = true
			}
		}
	case "persongroup":
		var persongroupCount int
		err = connection.DB.Self.Model(&user.UserGroup{}).
			Where("group in (?) and user = ?", currentStateValue["processor"].([]interface{}), c.GetInt("userId")).
			Count(&persongroupCount).Error
		if err != nil {
			return
		}
		if persongroupCount > 0 {
			status = true
		}
	case "department":
		var departmentCount int
		err = connection.DB.Self.Model(&user.Info{}).
			Where("dept in (?) and id = ?", currentStateValue["processor"].([]interface{}), c.GetInt("userId")).
			Count(&departmentCount).Error
		if err != nil {
			return
		}
		if departmentCount > 0 {
			status = true
		}
	case "variable":
		for _, p := range currentStateValue["processor"].([]interface{}) {
			switch int(p.(float64)) {
			case 1:
				if workOrderInfo.Creator == c.GetInt("userId") {
					status = true
				}
			case 2:
				err = connection.DB.Self.Model(&userInfo).Where("id = ?", workOrderInfo.Creator).Find(&userInfo).Error
				if err != nil {
					return
				}
				err = connection.DB.Self.Model(&userDept).Where("id = ?", userInfo.Dept).Find(&userDept).Error
				if err != nil {
					return
				}

				if userDept.Approver == c.GetInt("userId") {
					status = true
				} else if userDept.Leader == c.GetInt("userId") {
					status = true
				}
			}
		}
	}
	return
}
