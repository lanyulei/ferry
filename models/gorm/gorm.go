package gorm

import (
	"ferry/models/process"
	"ferry/models/system"

	"github.com/jinzhu/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		// 系统管理
		new(system.CasbinRule),
		new(system.Dept),
		new(system.Menu),
		new(system.LoginLog),
		new(system.RoleMenu),
		new(system.SysRoleDept),
		new(system.SysUser),
		new(system.SysRole),
		new(system.Post),
		new(system.Settings),

		// 流程中心
		new(process.Classify),
		new(process.TplInfo),
		new(process.TplData),
		new(process.WorkOrderInfo),
		new(process.TaskInfo),
		new(process.Info),
		new(process.History),
		new(process.CirculationHistory),
	).Error
}
