package service

import (
	"ferry/global/orm"
	"ferry/models/system"
	"strings"
)

/*
  @Author : lanyulei
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
