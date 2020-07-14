package gorm

import (
	"ferry/models/system"

	"github.com/jinzhu/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		// 系统管理
		&system.CasbinRule{},
		&system.Dept{},
		&system.Menu{},
		&system.LoginLog{},
		&system.RoleMenu{},
		&system.SysRoleDept{},
		&system.SysUser{},
		&system.SysRole{},
		&system.Post{},
		// 流程中心

	).Error
}
