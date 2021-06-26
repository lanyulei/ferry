package service

import (
	"encoding/json"
	"errors"
	"ferry/global/orm"
	"ferry/models/process"
	"ferry/models/system"
	"ferry/pkg/notify"
	"ferry/tools"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func CreateWorkOrder(c *gin.Context) (err error) {
	var (
		taskList       []string
		stateList      []interface{}
		userInfo       system.SysUser
		variableValue  []interface{}
		processValue   process.Info
		sendToUserList []system.SysUser
		noticeList     []int
		handle         Handle
		processState   ProcessState
		condExprStatus bool
		tpl            []byte
		sourceEdges    []map[string]interface{}
		targetEdges    []map[string]interface{}
		currentNode    map[string]interface{}
		workOrderValue struct {
			process.WorkOrderInfo
			Tpls        map[string][]interface{} `json:"tpls"`
			SourceState string                   `json:"source_state"`
			Tasks       json.RawMessage          `json:"tasks"`
			Source      string                   `json:"source"`
			IsExecTask  bool                     `json:"is_exec_task"`
		}
		paramsValue struct {
			Id       int           `json:"id"`
			Title    string        `json:"title"`
			Priority int           `json:"priority"`
			FormData []interface{} `json:"form_data"`
		}
	)

	err = c.ShouldBind(&workOrderValue)
	if err != nil {
		return
	}

	relatedPerson, err := json.Marshal([]int{tools.GetUserId(c)})
	if err != nil {
		return
	}

	// 获取变量值
	err = json.Unmarshal(workOrderValue.State, &variableValue)
	if err != nil {
		return
	}
	err = GetVariableValue(variableValue, tools.GetUserId(c))
	if err != nil {
		err = fmt.Errorf("获取处理人变量值失败，%v", err.Error())
		return
	}

	// 创建工单数据
	tx := orm.Eloquent.Begin()

	// 查询流程信息
	err = tx.Model(&processValue).Where("id = ?", workOrderValue.Process).Find(&processValue).Error
	if err != nil {
		return
	}

	err = json.Unmarshal(processValue.Structure, &processState.Structure)

	for _, node := range processState.Structure["nodes"] {
		if node["clazz"] == "start" {
			currentNode = node
		}
	}

	nodeValue, err := processState.GetNode(variableValue[0].(map[string]interface{})["id"].(string))
	if err != nil {
		return
	}

	for _, v := range workOrderValue.Tpls["form_data"] {
		tpl, err = json.Marshal(v)
		if err != nil {
			return
		}
		handle.WorkOrderData = append(handle.WorkOrderData, tpl)
	}

	switch nodeValue["clazz"] {
	// 排他网关
	case "exclusiveGateway":
		var sourceEdges []map[string]interface{}
		sourceEdges, err = processState.GetEdge(nodeValue["id"].(string), "source")
		if err != nil {
			return
		}
	breakTag:
		for _, edge := range sourceEdges {
			edgeCondExpr := make([]map[string]interface{}, 0)
			err = json.Unmarshal([]byte(edge["conditionExpression"].(string)), &edgeCondExpr)
			if err != nil {
				return
			}
			for _, condExpr := range edgeCondExpr {
				// 条件判断
				condExprStatus, err = handle.ConditionalJudgment(condExpr)
				if err != nil {
					return
				}
				if condExprStatus {
					// 进行节点跳转
					nodeValue, err = processState.GetNode(edge["target"].(string))
					if err != nil {
						return
					}

					if nodeValue["clazz"] == "userTask" || nodeValue["clazz"] == "receiveTask" {
						if nodeValue["assignValue"] == nil || nodeValue["assignType"] == "" {
							err = errors.New("处理人不能为空")
							return
						}
					}
					variableValue[0].(map[string]interface{})["id"] = nodeValue["id"].(string)
					variableValue[0].(map[string]interface{})["label"] = nodeValue["label"]
					variableValue[0].(map[string]interface{})["processor"] = nodeValue["assignValue"]
					variableValue[0].(map[string]interface{})["process_method"] = nodeValue["assignType"]
					break breakTag
				}
			}
		}
		if !condExprStatus {
			err = errors.New("所有流转均不符合条件，请确认。")
			return
		}
	case "parallelGateway":
		// 入口，判断
		sourceEdges, err = processState.GetEdge(nodeValue["id"].(string), "source")
		if err != nil {
			err = fmt.Errorf("查询流转信息失败，%v", err.Error())
			return
		}

		targetEdges, err = processState.GetEdge(nodeValue["id"].(string), "target")
		if err != nil {
			err = fmt.Errorf("查询流转信息失败，%v", err.Error())
			return
		}

		if len(sourceEdges) > 0 {
			nodeValue, err = processState.GetNode(sourceEdges[0]["target"].(string))
			if err != nil {
				return
			}
		} else {
			err = errors.New("并行网关流程不正确")
			return
		}

		if len(sourceEdges) > 1 && len(targetEdges) == 1 {
			// 入口
			variableValue = []interface{}{}
			for _, edge := range sourceEdges {
				targetStateValue, err := processState.GetNode(edge["target"].(string))
				if err != nil {
					return err
				}
				variableValue = append(variableValue, map[string]interface{}{
					"id":             edge["target"].(string),
					"label":          targetStateValue["label"],
					"processor":      targetStateValue["assignValue"],
					"process_method": targetStateValue["assignType"],
				})
			}
		} else {
			err = errors.New("并行网关流程配置不正确")
			return
		}
	}

	// 获取变量数据
	err = GetVariableValue(variableValue, tools.GetUserId(c))
	if err != nil {
		return
	}

	workOrderValue.State, err = json.Marshal(variableValue)
	if err != nil {
		return
	}

	var workOrderInfo = process.WorkOrderInfo{
		Title:         workOrderValue.Title,
		Priority:      workOrderValue.Priority,
		Process:       workOrderValue.Process,
		Classify:      workOrderValue.Classify,
		State:         workOrderValue.State,
		RelatedPerson: relatedPerson,
		Creator:       tools.GetUserId(c),
	}
	err = tx.Create(&workOrderInfo).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("创建工单失败，%v", err.Error())
		return
	}

	// 创建工单模版关联数据
	for i := 0; i < len(workOrderValue.Tpls["form_structure"]); i++ {
		var (
			formDataJson      []byte
			formStructureJson []byte
		)
		formDataJson, err = json.Marshal(workOrderValue.Tpls["form_data"][i])
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("生成json字符串错误，%v", err.Error())
			return
		}
		formStructureJson, err = json.Marshal(workOrderValue.Tpls["form_structure"][i])
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("生成json字符串错误，%v", err.Error())
			return
		}

		formData := process.TplData{
			WorkOrder:     workOrderInfo.Id,
			FormStructure: formStructureJson,
			FormData:      formDataJson,
		}

		err = tx.Create(&formData).Error
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("创建工单模版关联数据失败，%v", err.Error())
			return
		}
	}

	// 获取当前用户信息
	err = tx.Model(&system.SysUser{}).Where("user_id = ?", tools.GetUserId(c)).Find(&userInfo).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("查询用户信息失败，%v", err.Error())
		return
	}

	nameValue := userInfo.NickName
	if nameValue == "" {
		nameValue = userInfo.Username
	}

	// 创建历史记录
	err = json.Unmarshal(workOrderInfo.State, &stateList)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("json序列化失败，%s", err.Error())
		return
	}
	err = tx.Create(&process.CirculationHistory{
		Title:       workOrderValue.Title,
		WorkOrder:   workOrderInfo.Id,
		State:       workOrderValue.SourceState,
		Source:      workOrderValue.Source,
		Target:      stateList[0].(map[string]interface{})["id"].(string),
		Circulation: "新建",
		Processor:   nameValue,
		ProcessorId: userInfo.UserId,
		Status:      2, // 其他
	}).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("新建历史记录失败，%v", err.Error())
		return
	}

	// 更新流程提交数量统计
	err = tx.Model(&process.Info{}).
		Where("id = ?", workOrderValue.Process).
		Update("submit_count", processValue.SubmitCount+1).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("更新流程提交数量统计失败，%v", err.Error())
		return
	}

	tx.Commit()

	// 发送通知
	err = json.Unmarshal(processValue.Notice, &noticeList)
	if err != nil {
		return
	}
	if len(noticeList) > 0 {
		sendToUserList, err = GetPrincipalUserInfo(stateList, workOrderInfo.Creator)
		if err != nil {
			err = fmt.Errorf("获取所有处理人的用户信息失败，%v", err.Error())
			return
		}

		// 获取需要抄送的邮件
		emailCCList := make([]string, 0)
		if currentNode["cc"] != nil && len(currentNode["cc"].([]interface{})) > 0 {
			err = orm.Eloquent.Model(&system.SysUser{}).
				Where("user_id in (?)", currentNode["cc"]).
				Pluck("email", &emailCCList).Error
			if err != nil {
				err = errors.New("查询邮件抄送人失败")
				return
			}
		}
		// 发送通知
		go func() {
			bodyData := notify.BodyData{
				SendTo: map[string]interface{}{
					"userList": sendToUserList,
				},
				EmailCcTo:   emailCCList,
				Subject:     "您有一条待办工单，请及时处理",
				Description: "您有一条待办工单请及时处理，工单描述如下",
				Classify:    noticeList,
				ProcessId:   workOrderValue.Process,
				Id:          workOrderInfo.Id,
				Title:       workOrderValue.Title,
				Creator:     userInfo.NickName,
				Priority:    workOrderValue.Priority,
				CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			}
			err = bodyData.SendNotify()
			if err != nil {
				err = fmt.Errorf("通知发送失败，%v", err.Error())
				return
			}
		}()
	}

	if workOrderValue.IsExecTask {
		// 执行任务
		err = json.Unmarshal(workOrderValue.Tasks, &taskList)
		if err != nil {
			return
		}
		if len(taskList) > 0 {
			paramsValue.Id = workOrderInfo.Id
			paramsValue.Title = workOrderInfo.Title
			paramsValue.Priority = workOrderInfo.Priority
			paramsValue.FormData = workOrderValue.Tpls["form_data"]
			var params []byte
			params, err = json.Marshal(paramsValue)
			if err != nil {
				return
			}

			go ExecTask(taskList, string(params))
		}
	}

	return
}
