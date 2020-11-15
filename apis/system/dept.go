package system

import (
	"ferry/global/orm"
	"ferry/models/system"
	"ferry/tools"
	"ferry/tools/app"
	"ferry/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

/*
  @Author : lanyulei
*/

// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/deptList [get]
// @Security
func GetDeptList(c *gin.Context) {
	var (
		Dept   system.Dept
		err    error
		result []system.Dept
	)
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId, _ = tools.StringToInt(c.Request.FormValue("deptId"))

	if Dept.DeptName == "" {
		result, err = Dept.SetDept(true)
	} else {
		result, err = Dept.GetPage(true)
	}
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	app.OK(c, result, "")
}

func GetOrdinaryDeptList(c *gin.Context) {
	var (
		err      error
		deptList []system.Dept
	)

	err = orm.Eloquent.Model(&system.Dept{}).Find(&deptList).Error
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}

	app.OK(c, deptList, "")
}

func GetDeptTree(c *gin.Context) {
	var (
		Dept system.Dept
		err  error
	)
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId, _ = tools.StringToInt(c.Request.FormValue("deptId"))

	result, err := Dept.SetDept(false)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, "")
}

// @Summary 部门列表数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{deptId} [get]
// @Security
func GetDept(c *gin.Context) {
	var (
		err  error
		Dept system.Dept
	)
	Dept.DeptId, _ = tools.StringToInt(c.Param("deptId"))

	result, err := Dept.Get()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, msg.GetSuccess)
}

// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body models.Dept true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [post]
// @Security Bearer
func InsertDept(c *gin.Context) {
	var data system.Dept
	err := c.BindWith(&data, binding.JSON)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	data.CreateBy = tools.GetUserIdStr(c)
	result, err := data.Create()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, msg.CreatedSuccess)
}

// @Summary 修改部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body models.Dept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [put]
// @Security Bearer
func UpdateDept(c *gin.Context) {
	var data system.Dept
	err := c.BindJSON(&data)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.DeptId)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, msg.UpdatedSuccess)
}

// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dept/{id} [delete]
func DeleteDept(c *gin.Context) {
	var data system.Dept
	id, _ := tools.StringToInt(c.Param("id"))

	_, err := data.Delete(id)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, "", msg.DeletedSuccess)
}

func GetDeptTreeRoleSelect(c *gin.Context) {
	var Dept system.Dept
	var SysRole system.SysRole
	id, _ := tools.StringToInt(c.Param("roleId"))

	SysRole.RoleId = id
	result, err := Dept.SetDeptLable()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleDeptId()
		if err != nil {
			app.Error(c, -1, err, "")
			return
		}
	}
	app.Custum(c, gin.H{
		"code":        200,
		"depts":       result,
		"checkedKeys": menuIds,
	})
}
