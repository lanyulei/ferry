package process

import (
	"encoding/json"
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 工单
type WorkOrderInfo struct {
	base.Model
	Title         string          `gorm:"column:title; type:varchar(128)" json:"title" form:"title"`                      // 工单标题
	Process       int             `gorm:"column:process; type:int(11)" json:"process" form:"process"`                     // 流程ID
	Classify      int             `gorm:"column:classify; type:int(11)" json:"classify" form:"classify"`                  // 分类ID
	IsEnd         int             `gorm:"column:is_end; type:int(11); default:0" json:"is_end" form:"is_end"`             // 是否结束， 0 未结束，1 已结束
	State         json.RawMessage `gorm:"column:state; type:json" json:"state" form:"state"`                              // 状态信息
	RelatedPerson json.RawMessage `gorm:"column:related_person; type:json" json:"related_person" form:"related_person"`   // 工单所有处理人
	Creator       int             `gorm:"column:creator; type:int(11)" json:"creator" form:"creator"`                     // 创建人
	OrderType     int             `gorm:"column:order_type; type:int(11); default:1" json:"order_type" form:"order_type"` // 工单类型，1：正常流程，2：加签
}

func (WorkOrderInfo) TableName() string {
	return "p_work_order_info"
}
