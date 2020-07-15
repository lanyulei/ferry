package process

import (
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 任务
type History struct {
	base.Model
	Task          int    `gorm:"column:task; type: int(11)" json:"task" form:"task"`                                    // 任务ID
	Name          string `gorm:"column:name; type: varchar(256)" json:"name" form:"name"`                               // 任务名称
	TaskType      int    `gorm:"column:task_type; type: int(11)" json:"task_type" form:"task_type"`                     // 任务类型, python, shell
	ExecutionTime string `gorm:"column:execution_time; type: varchar(128)" json:"execution_time" form:"execution_time"` // 执行时间
	Result        string `gorm:"column:result; type: longtext" json:"result" form:"result"`                             // 任务返回
}

func (History) TableName() string {
	return "p_task_history"
}
