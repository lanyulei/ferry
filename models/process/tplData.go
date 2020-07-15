package process

import (
	"encoding/json"
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 工单绑定模版数据
type TplData struct {
	base.Model
	WorkOrder     int             `gorm:"column:work_order; type: int(11)" json:"work_order" form:"work_order"`          // 工单ID
	FormStructure json.RawMessage `gorm:"column:form_structure; type: json" json:"form_structure" form:"form_structure"` // 表单结构
	FormData      json.RawMessage `gorm:"column:form_data; type: json" json:"form_data" form:"form_data"`                // 表单数据
}

func (TplData) TableName() string {
	return "p_work_order_tpl_data"
}
