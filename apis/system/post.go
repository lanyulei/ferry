package system

import (
	"ferry/models/system"
	"ferry/tools"
	"ferry/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// @Summary 职位列表数据
// @Description 获取JSON
// @Tags 职位
// @Param postName query string false "postName"
// @Param postCode query string false "postCode"
// @Param postId query string false "postId"
// @Param status query string false "status"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post [get]
// @Security
func GetPostList(c *gin.Context) {
	var (
		data      system.Post
		err       error
		pageSize  = 10
		pageIndex = 1
	)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.PostId, _ = tools.StringToInt(c.Request.FormValue("postId"))

	data.PostCode = c.Request.FormValue("postCode")
	data.PostName = c.Request.FormValue("postName")
	data.Status = c.Request.FormValue("status")

	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

// @Summary 获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param postId path int true "postId"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post/{postId} [get]
// @Security
func GetPost(c *gin.Context) {
	var (
		err  error
		Post system.Post
	)
	Post.PostId, err = tools.StringToInt(c.Param("postId"))

	result, err := Post.Get()
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, "")
}

// @Summary 添加职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body models.Post true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post [post]
// @Security Bearer
func InsertPost(c *gin.Context) {
	var data system.Post
	err := c.Bind(&data)
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
	app.OK(c, result, "")
}

// @Summary 修改职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body models.Dept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post/ [put]
// @Security Bearer
func UpdatePost(c *gin.Context) {
	var data system.Post

	err := c.Bind(&data)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.PostId)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, "修改成功")
}

// @Summary 删除职位
// @Description 删除数据
// @Tags 职位
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/post/{postId} [delete]
func DeletePost(c *gin.Context) {
	var data system.Post
	data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup("postId", c)
	result, err := data.BatchDelete(IDS)
	if err != nil {
		app.Error(c, -1, err, "")
		return
	}
	app.OK(c, result, "删除成功")
}
