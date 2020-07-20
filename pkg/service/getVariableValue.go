package service

import (
	"ferry/global/orm"
	"ferry/models/system"
)

/*
  @Author : lanyulei
*/

func GetVariableValue(stateList []interface{}, creator int) (err error) {
	var (
		userInfo system.SysUser
		deptInfo system.Dept
	)

	// 变量转为实际的数据
	for _, stateItem := range stateList {
		if stateItem.(map[string]interface{})["process_method"] == "variable" {
			for processorIndex, processor := range stateItem.(map[string]interface{})["processor"].([]interface{}) {
				if int(processor.(float64)) == 1 {
					// 创建者
					stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = creator
				} else if int(processor.(float64)) == 2 {
					// 1. 查询用户信息
					err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
					if err != nil {
						return
					}
					// 2. 查询部门信息
					err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
					if err != nil {
						return
					}

					// 3. 替换处理人信息
					stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = deptInfo.Leader
				}
			}
			stateItem.(map[string]interface{})["process_method"] = "person"
		}
	}

	return
}
