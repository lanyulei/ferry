package process

import (
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 任务
type TaskInfo struct {
	base.Model
	Name     string `gorm:"column:name; type: varchar(256)" json:"name" form:"name"`               // 任务名称
	TaskType string `gorm:"column:task_type; type: varchar(45)" json:"task_type" form:"task_type"` // 任务类型
	Content  string `gorm:"column:content; type: longtext" json:"content" form:"content"`          // 任务内容
	Creator  int    `gorm:"column:creator; type: int(11)" json:"creator" form:"creator"`           // 创建者
	Remarks  string `gorm:"column:remarks; type: longtext" json:"remarks" form:"remarks"`          // 备注
}

func (TaskInfo) TableName() string {
	return "p_task_info"
}
