package system

import (
	"ferry/models/system"
	"ferry/tools/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// @Summary RoleMenu列表数据
// @Description 获取JSON
// @Tags 角色菜单
// @Param RoleId query string false "RoleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/rolemenu [get]
// @Security Bearer
func GetRoleMenu(c *gin.Context) {
	var (
		res app.Response
		Rm  system.RoleMenu
	)
	_ = c.ShouldBind(&Rm)
	result, err := Rm.Get()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

type RoleMenuPost struct {
	RoleId   string
	RoleMenu []system.RoleMenu
}

func InsertRoleMenu(c *gin.Context) {

	var res app.Response
	res.Msg = "添加成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 删除用户菜单数据
// @Description 删除数据
// @Tags 角色菜单
// @Param id path string true "id"
// @Param menu_id query string false "menu_id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/rolemenu/{id} [delete]
func DeleteRoleMenu(c *gin.Context) {
	var t system.RoleMenu
	id := c.Param("id")
	menuId := c.Request.FormValue("menu_id")
	_, err := t.Delete(id, menuId)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	var res app.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
