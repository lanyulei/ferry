package service

import (
	"errors"
	"ferry/global/orm"
	"ferry/models/system"
	"reflect"
	"strings"
)

/*
  @Author : lanyulei
  @todo: 添加新的处理人时候，需要修改（先完善功能，后续有时间的时候优化一下这部分。）
*/

func GetPrincipal(processor []int, processMethod string) (principals string, err error) {
	/*
		person 人员
		persongroup 人员组
		department 部门
		variable 变量
	*/
	var principalList []string
	switch processMethod {
	case "person":
		err = orm.Eloquent.Model(&system.SysUser{}).
			Where("user_id in (?)", processor).
			Pluck("nick_name", &principalList).Error
		if err != nil {
			return
		}
	case "role":
		err = orm.Eloquent.Model(&system.SysRole{}).
			Where("role_id in (?)", processor).
			Pluck("role_name", &principalList).Error
		if err != nil {
			return
		}
	case "department":
		err = orm.Eloquent.Model(&system.Dept{}).
			Where("dept_id in (?)", processor).
			Pluck("dept_name", &principalList).Error
		if err != nil {
			return
		}
	case "variable":
		for _, p := range processor {
			switch p {
			case 1:
				principalList = append(principalList, "创建者")
			case 2:
				principalList = append(principalList, "创建者负责人")
			}
		}
	}
	return strings.Join(principalList, ","), nil
}

// 获取用户对应
func GetPrincipalUserInfo(stateList []interface{}, creator int) (userInfoList []system.SysUser, err error) {
	var (
		userInfo        system.SysUser
		deptInfo        system.Dept
		userInfoListTmp []system.SysUser // 临时保存查询的列表数据
		processorList   []interface{}
	)

	err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
	if err != nil {
		return
	}

	for _, stateItem := range stateList {

		if reflect.TypeOf(stateItem.(map[string]interface{})["processor"]) == nil {
			err = errors.New("未找到对应的处理人，请确认。")
			return
		}
		stateItemType := reflect.TypeOf(stateItem.(map[string]interface{})["processor"]).String()
		if stateItemType == "[]int" {
			for _, v := range stateItem.(map[string]interface{})["processor"].([]int) {
				processorList = append(processorList, v)
			}
		} else {
			processorList = stateItem.(map[string]interface{})["processor"].([]interface{})
		}

		switch stateItem.(map[string]interface{})["process_method"] {
		case "person":
			err = orm.Eloquent.Model(&system.SysUser{}).
				Where("user_id in (?)", processorList).
				Find(&userInfoListTmp).Error
			if err != nil {
				return
			}
			userInfoList = append(userInfoList, userInfoListTmp...)
		case "role":
			err = orm.Eloquent.Model(&system.SysUser{}).
				Where("role_id in (?)", processorList).
				Find(&userInfoListTmp).Error
			if err != nil {
				return
			}
			userInfoList = append(userInfoList, userInfoListTmp...)
		case "department":
			err = orm.Eloquent.Model(&system.SysUser{}).
				Where("dept_id in (?)", processorList).
				Find(&userInfoListTmp).Error
			if err != nil {
				return
			}
			userInfoList = append(userInfoList, userInfoListTmp...)
		case "variable": // 变量
			for _, processor := range processorList {
				if int(processor.(float64)) == 1 {
					// 创建者
					userInfoList = append(userInfoList, userInfo)
				} else if int(processor.(float64)) == 2 {
					// 1. 查询部门信息
					err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
					if err != nil {
						return
					}

					// 2. 查询Leader信息
					err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", deptInfo.Leader).Find(&userInfo).Error
					if err != nil {
						return
					}
					userInfoList = append(userInfoList, userInfo)
				}
			}
		}
	}

	return
}
