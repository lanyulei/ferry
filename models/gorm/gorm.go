package gorm

import (
	"ferry/models"
	"ferry/models/tools"

	"github.com/jinzhu/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		new(models.CasbinRule),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(models.Dept),
		new(models.Menu),
		new(models.LoginLog),
		new(models.RoleMenu),
		new(models.SysRoleDept),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Post),
	).Error
}
