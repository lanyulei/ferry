/*
  @Author : lanyulei
*/

package base

import (
	"ferry/pkg/jsonTime"
)

type Model struct {
	Id        int                `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id" form:"id"`
	CreatedAt jsonTime.JSONTime  `gorm:"column:create_time" json:"create_time" form:"create_time"`
	UpdatedAt jsonTime.JSONTime  `gorm:"column:update_time" json:"update_time" form:"update_time"`
	DeletedAt *jsonTime.JSONTime `gorm:"column:delete_time" sql:"index" json:"-"`
}
