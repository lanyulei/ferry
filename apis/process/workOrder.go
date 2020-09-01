package process

import (
	"encoding/json"
	"errors"
	"ferry/global/orm"
	"ferry/models/process"
	"ferry/models/system"
	"ferry/pkg/notify"
	"ferry/pkg/service"
	"ferry/tools"
	"ferry/tools/app"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
 @Author : lanyulei
*/

// 流程结构包括节点，流转和模版
func ProcessStructure(c *gin.Context) {
	processId := c.DefaultQuery("processId", "")
	if processId == "" {
		app.Error(c, -1, errors.New("参数不正确，请确定参数processId是否传递"), "")
		return
	}
	workOrderId := c.DefaultQuery("workOrderId", "0")
	if workOrderId == "" {
		app.Error(c, -1, errors.New("参数不正确，请确定参数workOrderId是否传递"), "")
		return
	}
	workOrderIdInt, _ := strconv.Atoi(workOrderId)
	processIdInt, _ := strconv.Atoi(processId)
	result, err := service.ProcessStructure(c, processIdInt, workOrderIdInt)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	if workOrderIdInt != 0 {
		currentState := result["workOrder"].(service.WorkOrderData).CurrentState
		userAuthority, err := service.JudgeUserAuthority(c, workOrderIdInt, currentState)
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("判断用户是否有权限失败，%v", err.Error()))
			return
		}
		result["userAuthority"] = userAuthority
	}

	app.OK(c, result, "数据获取成功")
}

// 新建工单
func CreateWorkOrder(c *gin.Context) {
	var (
		taskList       []string
		stateList      []interface{}
		userInfo       system.SysUser
		variableValue  []interface{}
		processValue   process.Info
		sendToUserList []system.SysUser
		noticeList     []int
		workOrderValue struct {
			process.WorkOrderInfo
			Tpls        map[string][]interface{} `json:"tpls"`
			SourceState string                   `json:"source_state"`
			Tasks       json.RawMessage          `json:"tasks"`
			Source      string                   `json:"source"`
		}
	)

	err := c.ShouldBind(&workOrderValue)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	relatedPerson, err := json.Marshal([]int{tools.GetUserId(c)})
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 获取变量值
	err = json.Unmarshal(workOrderValue.State, &variableValue)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	err = service.GetVariableValue(variableValue, tools.GetUserId(c))
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("获取处理人变量值失败，%v", err.Error()))
		return
	}
	workOrderValue.State, err = json.Marshal(variableValue)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 创建工单数据
	tx := orm.Eloquent.Begin()

	// 查询流程信息
	err = tx.Model(&processValue).Where("id = ?", workOrderValue.Process).Find(&processValue).Error
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
		app.Error(c, -1, err, fmt.Sprintf("创建工单失败，%v", err.Error()))
		return
	}

	// 创建工单模版关联数据
	for i := 0; i < len(workOrderValue.Tpls["form_structure"]); i++ {
		formDataJson, err := json.Marshal(workOrderValue.Tpls["form_data"][i])
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, fmt.Sprintf("生成json字符串错误，%v", err.Error()))
			return
		}
		formStructureJson, err := json.Marshal(workOrderValue.Tpls["form_structure"][i])
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, fmt.Sprintf("生成json字符串错误，%v", err.Error()))
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
			app.Error(c, -1, err, fmt.Sprintf("创建工单模版关联数据失败，%v", err.Error()))
			return
		}
	}

	// 获取当前用户信息
	err = tx.Model(&system.SysUser{}).Where("user_id = ?", tools.GetUserId(c)).Find(&userInfo).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, fmt.Sprintf("查询用户信息失败，%v", err.Error()))
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
		app.Error(c, -1, err, fmt.Sprintf("Json序列化失败，%v", err.Error()))
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
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, fmt.Sprintf("新建历史记录失败，%v", err.Error()))
		return
	}

	// 更新流程提交数量统计
	err = tx.Model(&process.Info{}).
		Where("id = ?", workOrderValue.Process).
		Update("submit_count", processValue.SubmitCount+1).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, fmt.Sprintf("更新流程提交数量统计失败，%v", err.Error()))
		return
	}

	tx.Commit()

	// 发送通知
	err = json.Unmarshal(processValue.Notice, &noticeList)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	if len(noticeList) > 0 {
		sendToUserList, err = service.GetPrincipalUserInfo(stateList, workOrderInfo.Creator)
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("获取所有处理人的用户信息失败，%v", err.Error()))
			return
		}

		// 发送通知
		go func() {
			bodyData := notify.BodyData{
				SendTo: map[string]interface{}{
					"userList": sendToUserList,
				},
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
				app.Error(c, -1, err, fmt.Sprintf("通知发送失败，%v", err.Error()))
				return
			}
		}()
	}

	// 执行任务
	err = json.Unmarshal(workOrderValue.Tasks, &taskList)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	if len(taskList) > 0 {
		go service.ExecTask(taskList)
	}

	app.OK(c, "", "成功提交工单申请")
}

// 工单列表
func WorkOrderList(c *gin.Context) {
	/*
		1. 待办工单
		2. 我创建的
		3. 我相关的
		4. 所有工单
	*/

	var (
		result      interface{}
		err         error
		classifyInt int
	)

	classify := c.DefaultQuery("classify", "0")
	if classify == "" {
		app.Error(c, -1, errors.New("参数错误，请确认classify是否传递"), "")
		return
	}

	classifyInt, _ = strconv.Atoi(classify)
	w := service.WorkOrder{
		Classify: classifyInt,
		GinObj:   c,
	}
	result, err = w.WorkOrderList()
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询工单数据失败，%v", err.Error()))
		return
	}

	app.OK(c, result, "")
}

// 处理工单
func ProcessWorkOrder(c *gin.Context) {
	var (
		err           error
		userAuthority bool
		handle        service.Handle
		params        struct {
			Tasks          []string
			TargetState    string                   `json:"target_state"`    // 目标状态
			SourceState    string                   `json:"source_state"`    // 源状态
			WorkOrderId    int                      `json:"work_order_id"`   // 工单ID
			Circulation    string                   `json:"circulation"`     // 流转ID
			FlowProperties int                      `json:"flow_properties"` // 流转类型 0 拒绝，1 同意，2 其他
			Remarks        string                   `json:"remarks"`         // 处理的备注信息
			Tpls           []map[string]interface{} `json:"tpls"`            // 表单数据
		}
	)

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 处理工单
	userAuthority, err = service.JudgeUserAuthority(c, params.WorkOrderId, params.SourceState)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("判断用户是否有权限失败，%v", err.Error()))
		return
	}
	if !userAuthority {
		app.Error(c, -1, errors.New("当前用户没有权限进行此操作"), "")
		return
	}

	err = handle.HandleWorkOrder(
		c,
		params.WorkOrderId,    // 工单ID
		params.Tasks,          // 任务列表
		params.TargetState,    // 目标节点
		params.SourceState,    // 源节点
		params.Circulation,    // 流转标题
		params.FlowProperties, // 流转属性
		params.Remarks,        // 备注信息
		params.Tpls,           // 工单数据更新
	)
	if err != nil {
		app.Error(c, -1, nil, fmt.Sprintf("处理工单失败，%v", err.Error()))
		return
	}

	app.OK(c, nil, "工单处理完成")
}

// 结束工单
func UnityWorkOrder(c *gin.Context) {
	var (
		err           error
		workOrderId   string
		workOrderInfo process.WorkOrderInfo
		userInfo      system.SysUser
	)

	workOrderId = c.DefaultQuery("work_oroder_id", "")
	if workOrderId == "" {
		app.Error(c, -1, errors.New("参数不正确，work_oroder_id"), "")
		return
	}

	tx := orm.Eloquent.Begin()

	// 查询工单信息
	err = tx.Model(&workOrderInfo).
		Where("id = ?", workOrderId).
		Find(&workOrderInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询工单失败，%v", err.Error()))
		return
	}
	if workOrderInfo.IsEnd == 1 {
		app.Error(c, -1, errors.New("工单已结束"), "")
		return
	}

	// 更新工单状态
	err = tx.Model(&process.WorkOrderInfo{}).
		Where("id = ?", workOrderId).
		Update("is_end", 1).
		Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, fmt.Sprintf("结束工单失败，%v", err.Error()))
		return
	}

	// 获取当前用户信息
	err = tx.Model(&userInfo).
		Where("user_id = ?", tools.GetUserId(c)).
		Find(&userInfo).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, fmt.Sprintf("当前用户查询失败，%v", err.Error()))
		return
	}

	// 写入历史
	tx.Create(&process.CirculationHistory{
		Title:       workOrderInfo.Title,
		WorkOrder:   workOrderInfo.Id,
		State:       "结束工单",
		Circulation: "结束",
		Processor:   userInfo.NickName,
		ProcessorId: tools.GetUserId(c),
		Remarks:     "手动结束工单。",
	})

	tx.Commit()

	app.OK(c, nil, "工单已结束")
}

// 转交工单
func InversionWorkOrder(c *gin.Context) {
	var (
		err             error
		workOrderInfo   process.WorkOrderInfo
		stateList       []map[string]interface{}
		stateValue      []byte
		currentState    map[string]interface{}
		userInfo        system.SysUser
		currentUserInfo system.SysUser
		params          struct {
			WorkOrderId int    `json:"work_order_id"`
			NodeId      string `json:"node_id"`
			UserId      int    `json:"user_id"`
			Remarks     string `json:"remarks"`
		}
	)

	// 获取当前用户信息
	err = orm.Eloquent.Model(&currentUserInfo).
		Where("user_id = ?", tools.GetUserId(c)).
		Find(&currentUserInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("当前用户查询失败，%v", err.Error()))
		return
	}

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	// 查询工单信息
	err = orm.Eloquent.Model(&workOrderInfo).
		Where("id = ?", params.WorkOrderId).
		Find(&workOrderInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询工单信息失败，%v", err.Error()))
		return
	}

	// 序列化节点数据
	err = json.Unmarshal(workOrderInfo.State, &stateList)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("节点数据反序列化失败，%v", err.Error()))
		return
	}

	for _, s := range stateList {
		if s["id"].(string) == params.NodeId {
			s["processor"] = []interface{}{params.UserId}
			s["process_method"] = "person"
			currentState = s
			break
		}
	}

	stateValue, err = json.Marshal(stateList)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("节点数据序列化失败，%v", err.Error()))
		return
	}

	tx := orm.Eloquent.Begin()

	// 更新数据
	err = tx.Model(&process.WorkOrderInfo{}).
		Where("id = ?", params.WorkOrderId).
		Update("state", stateValue).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("更新节点信息失败，%v", err.Error()))
		return
	}

	// 查询用户信息
	err = tx.Model(&system.SysUser{}).
		Where("user_id = ?", params.UserId).
		Find(&userInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询用户信息失败，%v", err.Error()))
		return
	}

	// 添加转交历史
	tx.Create(&process.CirculationHistory{
		Title:       workOrderInfo.Title,
		WorkOrder:   workOrderInfo.Id,
		State:       currentState["label"].(string),
		Circulation: "转交",
		Processor:   currentUserInfo.NickName,
		ProcessorId: tools.GetUserId(c),
		Remarks:     fmt.Sprintf("此阶段负责人已转交给《%v》", userInfo.NickName),
	})

	tx.Commit()

	app.OK(c, nil, "工单已手动结单")
}

// 催办工单
func UrgeWorkOrder(c *gin.Context) {
	var (
		workOrderInfo  process.WorkOrderInfo
		sendToUserList []system.SysUser
		stateList      []interface{}
		userInfo       system.SysUser
	)
	workOrderId := c.DefaultQuery("workOrderId", "")
	if workOrderId == "" {
		app.Error(c, -1, errors.New("参数不正确，缺失workOrderId"), "")
		return
	}

	// 查询工单数据
	err := orm.Eloquent.Model(&process.WorkOrderInfo{}).Where("id = ?", workOrderId).Find(&workOrderInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("查询工单信息失败，%v", err.Error()))
		return
	}

	// 确认是否可以催办
	if workOrderInfo.UrgeLastTime != 0 && (int(time.Now().Unix())-workOrderInfo.UrgeLastTime) < 600 {
		app.Error(c, -1, errors.New("十分钟内无法多次催办工单。"), "")
		return
	}

	// 获取当前工单处理人信息
	err = json.Unmarshal(workOrderInfo.State, &stateList)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	sendToUserList, err = service.GetPrincipalUserInfo(stateList, workOrderInfo.Creator)

	// 查询创建人信息
	err = orm.Eloquent.Model(&system.SysUser{}).Where("user_id = ?", workOrderInfo.Creator).Find(&userInfo).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("创建人信息查询失败，%v", err.Error()))
		return
	}

	// 发送催办提醒
	bodyData := notify.BodyData{
		SendTo: map[string]interface{}{
			"userList": sendToUserList,
		},
		Subject:     "您被催办工单了，请及时处理。",
		Description: "您有一条待办工单，请及时处理，工单描述如下",
		Classify:    []int{1}, // todo 1 表示邮箱，后续添加了其他的在重新补充
		ProcessId:   workOrderInfo.Process,
		Id:          workOrderInfo.Id,
		Title:       workOrderInfo.Title,
		Creator:     userInfo.NickName,
		Priority:    workOrderInfo.Priority,
		CreatedAt:   workOrderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	err = bodyData.SendNotify()
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("催办提醒发送失败，%v", err.Error()))
		return
	}

	// 更新数据库
	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).Where("id = ?", workOrderInfo.Id).Updates(map[string]interface{}{
		"urge_count":     workOrderInfo.UrgeCount + 1,
		"urge_last_time": int(time.Now().Unix()),
	}).Error
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("更新催办信息失败，%v", err.Error()))
		return
	}

	app.OK(c, "", "")
}

// 主动处理
func ActiveOrder(c *gin.Context) {
	var (
		workOrderId string
		err         error
		stateValue  []struct {
			ID            string `json:"id"`
			Label         string `json:"label"`
			ProcessMethod string `json:"process_method"`
			Processor     []int  `json:"processor"`
		}
		stateValueByte []byte
	)

	workOrderId = c.Param("id")

	err = c.ShouldBind(&stateValue)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	stateValueByte, err = json.Marshal(stateValue)
	if err != nil {
		app.Error(c, -1, fmt.Errorf("转byte失败，%v", err.Error()), "")
		return
	}

	err = orm.Eloquent.Model(&process.WorkOrderInfo{}).
		Where("id = ?", workOrderId).
		Update("state", stateValueByte).Error
	if err != nil {
		app.Error(c, -1, fmt.Errorf("接单失败，%v", err.Error()), "")
		return
	}

	app.OK(c, "", "接单成功，请及时处理")
}
