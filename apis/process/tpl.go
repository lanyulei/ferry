package process

//import (
//	"ferry/models/tpl"
//	"ferry/pkg/connection"
//	"ferry/pkg/pagination"
//	"ferry/pkg/response/code"
//	. "ferry/pkg/response/response"
//	"fmt"
//
//	"github.com/gin-gonic/gin"
//)
//
///*
//  @Author : lanyulei
//*/
//
//// 模板列表
//func TemplateList(c *gin.Context) {
//	type templateUserValue struct {
//		tpl.Info
//		CreateUser string `json:"create_user"`
//		CreateName string `json:"create_name"`
//	}
//
//	var (
//		err          error
//		templateList []*templateUserValue
//	)
//
//	SearchParams := map[string]map[string]interface{}{
//		"like": pagination.RequestParams(c),
//	}
//
//	db := connection.DB.Self.Model(&tpl.Info{}).Joins("left join user_info on user_info.id = tpl_info.creator").
//		Select("tpl_info.id, tpl_info.create_time, tpl_info.update_time, tpl_info.`name`, tpl_info.`creator`, user_info.username as create_user, user_info.nickname as create_name").Where("tpl_info.`delete_time` IS NULL")
//
//	result, err := pagination.Paging(&pagination.Param{
//		C:  c,
//		DB: db,
//	}, &templateList, SearchParams, "tpl_info")
//
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询模版失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, result, "")
//}
//
//// 创建模版
//func CreateTemplate(c *gin.Context) {
//	var (
//		err           error
//		templateValue tpl.Info
//		templateCount int
//	)
//
//	err = c.ShouldBind(&templateValue)
//	if err != nil {
//		Response(c, code.BindError, nil, err.Error())
//		return
//	}
//
//	// 确定修改的分类是否存在
//	err = connection.DB.Self.Model(&templateValue).
//		Where("name = ?", templateValue.Name).
//		Count(&templateCount).Error
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询模版数量失败，%v", err.Error()))
//		return
//	}
//	if templateCount > 0 {
//		Response(c, code.InternalServerError, nil, "模版名称出现重复，请换一个名称")
//		return
//	}
//
//	templateValue.Creator = c.GetInt("userId") // 当前登陆用户ID
//	err = connection.DB.Self.Create(&templateValue).Error
//	if err != nil {
//		Response(c, code.CreateError, nil, fmt.Sprintf("创建模板失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 模版详情
//func TemplateDetails(c *gin.Context) {
//	var (
//		err                  error
//		templateDetailsValue tpl.Info
//	)
//
//	templateId := c.DefaultQuery("template_id", "")
//	if templateId == "" {
//		Response(c, code.ParamError, nil, fmt.Sprintf("参数不正确，请确认template_id是否传递"))
//		return
//	}
//
//	err = connection.DB.Self.Model(&templateDetailsValue).Where("id = ?", templateId).Find(&templateDetailsValue).Error
//	if err != nil {
//		Response(c, code.SelectError, nil, fmt.Sprintf("查询模版数据失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, templateDetailsValue, "")
//}
//
//// 更新模版
//func UpdateTemplate(c *gin.Context) {
//	var (
//		err           error
//		templateValue tpl.Info
//	)
//	err = c.ShouldBind(&templateValue)
//	if err != nil {
//		Response(c, code.BindError, nil, fmt.Sprintf("参数绑定失败，%v", err.Error()))
//		return
//	}
//
//	err = connection.DB.Self.Model(&templateValue).Where("id = ?", templateValue.Id).Updates(map[string]interface{}{
//		"name":           templateValue.Name,
//		"remarks":        templateValue.Remarks,
//		"form_structure": templateValue.FormStructure,
//	}).Error
//	if err != nil {
//		Response(c, code.UpdateError, nil, fmt.Sprintf("更新模版失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, templateValue, "")
//}
//
//// 删除模版
//func DeleteTemplate(c *gin.Context) {
//	var (
//		err error
//	)
//
//	templateId := c.DefaultQuery("template_id", "")
//	if templateId == "" {
//		Response(c, code.ParamError, nil, fmt.Sprintf("参数不正确，请确认template_id是否传递"))
//		return
//	}
//
//	err = connection.DB.Self.Delete(tpl.Info{}, "id = ?", templateId).Error
//	if err != nil {
//		Response(c, code.DeleteError, nil, fmt.Sprintf("模版删除失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
