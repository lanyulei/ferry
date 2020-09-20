package process

import (
	"encoding/json"
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 流程
type Info struct {
	base.Model
	Name        string          `gorm:"column:name; type:varchar(128)" json:"name" form:"name"`                               // 流程名称
	Icon        string          `gorm:"column:icon; type:varchar(128)" json:"icon" form:"icon"`                               // 流程标签
	Structure   json.RawMessage `gorm:"column:structure; type:json" json:"structure" form:"structure"`                        // 流程结构
	Classify    int             `gorm:"column:classify; type:int(11)" json:"classify" form:"classify"`                        // 分类ID
	Tpls        json.RawMessage `gorm:"column:tpls; type:json" json:"tpls" form:"tpls"`                                       // 模版
	Task        json.RawMessage `gorm:"column:task; type:json" json:"task" form:"task"`                                       // 任务ID, array, 可执行多个任务，可以当成通知任务，每个节点都会去执行
	SubmitCount int             `gorm:"column:submit_count; type:int(11); default:0" json:"submit_count" form:"submit_count"` // 提交统计
	Creator     int             `gorm:"column:creator; type:int(11)" json:"creator" form:"creator"`                           // 创建者
	Notice      json.RawMessage `gorm:"column:notice; type:json" json:"notice" form:"notice"`                                 // 绑定通知
	Remarks     string          `gorm:"column:remarks; type:varchar(1024)" json:"remarks" form:"remarks"`                     // 流程备注
}

func (Info) TableName() string {
	return "p_process_info"
}
