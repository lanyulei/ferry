package system

import (
	"ferry/global/orm"
	"ferry/models/system"
	"ferry/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 设置系统信息
func GetSettingsInfo(c *gin.Context) {
	var (
		err          error
		settingsInfo []*system.Settings
		classify     string
	)
	db := orm.Eloquent.Model(&settingsInfo)
	classify = c.DefaultQuery("classify", "")
	if classify != "" {
		db = db.Where("classify = ?", classify)
	}

	err = db.Find(&settingsInfo).Error
	if err != nil {
		app.Error(c, -1, fmt.Errorf("查询数据失败，%v", err.Error()), "")
		return
	}

	app.OK(c, settingsInfo, "查询配置信息成功")
}

// 设置系统信息
func SetSettingsInfo(c *gin.Context) {
	var (
		err           error
		settingsInfo  system.Settings
		settingsCount int
	)

	err = c.ShouldBind(&settingsInfo)
	if err != nil {
		app.Error(c, -1, fmt.Errorf("绑定数据失败，%v", err.Error()), "")
		return
	}

	// 查询数据是否存在
	err = orm.Eloquent.Model(&system.Settings{}).
		Where("classify = ?", settingsInfo.Classify).
		Count(&settingsCount).Error
	if err != nil {
		app.Error(c, -1, fmt.Errorf("查询数据失败，%v", err.Error()), "")
		return
	}
	if settingsCount == 0 {
		// 创建新的配置信息
		err = orm.Eloquent.Create(&settingsInfo).Error
		if err != nil {
			app.Error(c, -1, fmt.Errorf("创建配置信息失败，%v", err.Error()), "")
			return
		}
	} else {
		err = orm.Eloquent.Model(&settingsInfo).
			Where("classify = ?", settingsInfo.Classify).
			Updates(&settingsInfo).Error
		if err != nil {
			app.Error(c, -1, fmt.Errorf("更新配置信息失败，%v", err.Error()), "")
			return
		}
	}

	app.OK(c, "", "配置信息设置成功")
}
