package system

import (
	"ferry/models/system"
	"ferry/tools"
	"ferry/tools/app"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

/*
  @Author : lanyulei
*/

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security
func GetRoleList(c *gin.Context) {
	var (
		err       error
		pageSize  = 10
		pageIndex = 1
		data      system.SysRole
	)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role [get]
// @Security Bearer
func GetRole(c *gin.Context) {
	var (
		err  error
		Role system.SysRole
	)
	Role.RoleId, _ = tools.StringToInt(c.Param("roleId"))

	result, err := Role.Get()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	menuIds, err := Role.GetRoleMeunId()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	result.MenuIds = menuIds
	app.OK(c, result, "")

}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
func InsertRole(c *gin.Context) {
	var data system.SysRole
	data.CreateBy = tools.GetUserIdStr(c)
	err := c.BindWith(&data, binding.JSON)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	id, err := data.Insert()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	data.RoleId = id
	var t system.RoleMenu
	_, err = t.Insert(id, data.MenuIds)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, data, "添加成功")
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role [put]
func UpdateRole(c *gin.Context) {
	var (
		data system.SysRole
		t    system.RoleMenu
		err  error
	)
	data.UpdateBy = tools.GetUserIdStr(c)
	err = c.Bind(&data)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	result, err := data.Update(data.RoleId)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	_, err = t.DeleteRoleMenu(data.RoleId)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	_, err = t.Insert(data.RoleId, data.MenuIds)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	app.OK(c, result, "修改成功")
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param roleId path int true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role/{roleId} [delete]
func DeleteRole(c *gin.Context) {
	var Role system.SysRole
	Role.UpdateBy = tools.GetUserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup("roleId", c)
	_, err := Role.BatchDelete(IDS)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	var t system.RoleMenu
	_, err = t.BatchDeleteRoleMenu(IDS)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, "", "删除成功")
}
