package process

import (
	"encoding/json"
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

// 模板
type TplInfo struct {
	base.Model
	Name          string          `gorm:"column:name; type: varchar(128)" json:"name" form:"name" binding:"required"`                       // 模板名称
	FormStructure json.RawMessage `gorm:"column:form_structure; type: json" json:"form_structure" form:"form_structure" binding:"required"` // 表单结构
	Creator       int             `gorm:"column:creator; type: int(11)" json:"creator" form:"creator"`                                      // 创建者
	Remarks       string          `gorm:"column:remarks; type: longtext" json:"remarks" form:"remarks"`                                     // 备注
}

func (TplInfo) TableName() string {
	return "p_tpl_info"
}
