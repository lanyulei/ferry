package system

import (
	"ferry/global/orm"
	"ferry/tools"
)

/*
  @Author : lanyulei
*/

type Login struct {
	Username  string `form:"UserName" json:"username" binding:"required"`
	Password  string `form:"Password" json:"password" binding:"required"`
	Code      string `form:"Code" json:"code" binding:"required"`
	UUID      string `form:"UUID" json:"uuid" binding:"required"`
	LoginType int    `form:"LoginType" json:"loginType"`
}

func (u *Login) GetUser() (user SysUser, role SysRole, e error) {

	e = orm.Eloquent.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
	if e != nil {
		return
	}

	// 验证密码
	if u.LoginType == 0 {
		_, e = tools.CompareHashAndPassword(user.Password, u.Password)
		if e != nil {
			return
		}
	}

	e = orm.Eloquent.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		return
	}
	return
}
