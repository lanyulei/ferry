package process

import (
	"errors"
	"ferry/global/orm"
	"ferry/models/process"
	"ferry/pkg/pagination"
	"ferry/tools"
	"ferry/tools/app"

	"github.com/gin-gonic/gin"
)

/*
 @Author : lanyulei
*/

// 模板列表
func TemplateList(c *gin.Context) {
	type templateUserValue struct {
		process.TplInfo
		CreateUser string `json:"create_user"`
		CreateName string `json:"create_name"`
	}

	var (
		err          error
		templateList []*templateUserValue
	)

	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	db := orm.Eloquent.Model(&process.TplInfo{}).Joins("left join sys_user on sys_user.user_id = p_tpl_info.creator").
		Select("p_tpl_info.id, p_tpl_info.create_time, p_tpl_info.update_time, p_tpl_info.`name`, p_tpl_info.`creator`, sys_user.username as create_user, sys_user.nick_name as create_name")

	result, err := pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &templateList, SearchParams, "p_tpl_info")

	if err != nil {
		tools.HasError(err, "")
	}

	app.OK(c, result, "获取模版列表成功")
}

// 创建模版
func CreateTemplate(c *gin.Context) {
	var (
		err           error
		templateValue process.TplInfo
		templateCount int
	)

	err = c.ShouldBind(&templateValue)
	if err != nil {
		tools.HasError(err, "")
	}

	// 确认模版名称是否存在
	err = orm.Eloquent.Model(&templateValue).
		Where("name = ?", templateValue.Name).
		Count(&templateCount).Error
	if err != nil {
		tools.HasError(err, "")
	}
	if templateCount > 0 {
		tools.HasError(errors.New("模版名称已存在"), "")
	}

	templateValue.Creator = tools.GetUserId(c) // 当前登陆用户ID
	err = orm.Eloquent.Create(&templateValue).Error
	if err != nil {
		tools.HasError(err, "")
	}

	app.OK(c, "", "创建模版成功")
}

// 模版详情
func TemplateDetails(c *gin.Context) {
	var (
		err                  error
		templateDetailsValue process.TplInfo
	)

	templateId := c.DefaultQuery("template_id", "")
	if templateId == "" {
		tools.HasError(err, "")
	}

	err = orm.Eloquent.Model(&templateDetailsValue).
		Where("id = ?", templateId).
		Find(&templateDetailsValue).Error
	if err != nil {
		tools.HasError(err, "")
	}

	app.OK(c, templateDetailsValue, "")
}

// 更新模版
func UpdateTemplate(c *gin.Context) {
	var (
		err           error
		templateValue process.TplInfo
	)
	err = c.ShouldBind(&templateValue)
	if err != nil {
		tools.HasError(err, "")
	}

	err = orm.Eloquent.Model(&templateValue).
		Where("id = ?", templateValue.Id).
		Updates(map[string]interface{}{
			"name":           templateValue.Name,
			"remarks":        templateValue.Remarks,
			"form_structure": templateValue.FormStructure,
		}).Error
	if err != nil {
		tools.HasError(err, "")
	}

	app.OK(c, templateValue, "")
}

// 删除模版
func DeleteTemplate(c *gin.Context) {
	var (
		err error
	)

	templateId := c.DefaultQuery("templateId", "")
	if templateId == "" {
		tools.HasError(errors.New("参数不正确，请确认templateId是否传递"), "")
	}

	err = orm.Eloquent.Delete(process.TplInfo{}, "id = ?", templateId).Error
	if err != nil {
		tools.HasError(err, "")
	}

	app.OK(c, "", "删除模版成功")
}
