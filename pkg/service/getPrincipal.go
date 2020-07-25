package service

import (
	"ferry/global/orm"
	"ferry/models/system"
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
	//case "department":
	//	err = orm.Eloquent.Model(&user.Dept{}).Where("id in (?)", processor).Pluck("nickname", &principalList).Error
	//	if err != nil {
	//		return
	//	}
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
func GetPrincipalUserInfo(stateList []map[string]interface{}, creator int) (userInfoList []system.SysUser, err error) {
	var (
		userInfo        system.SysUser
		deptInfo        system.Dept
		userInfoListTmp []system.SysUser // 临时保存查询的列表数据
	)

	err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
	if err != nil {
		return
	}

	for _, stateItem := range stateList {
		switch stateItem["process_method"] {
		case "person":
			err = orm.Eloquent.Model(&system.SysUser{}).
				Where("user_id in (?)", stateItem["processor"].([]interface{})).
				Find(&userInfoListTmp).Error
			if err != nil {
				return
			}
			userInfoList = append(userInfoList, userInfoListTmp...)
		case "variable": // 变量
			for _, processor := range stateItem["processor"].([]interface{}) {
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
