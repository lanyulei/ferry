package public

import (
	"encoding/base64"
	"errors"
	"ferry/pkg/utils"
	"ferry/tools/app"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary 上传图片
// @Description 获取JSON
// @Tags 公共接口
// @Accept multipart/form-data
// @Param type query string true "type" (1：单图，2：多图, 3：base64图片)
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/public/uploadFile [post]

func UploadFile(c *gin.Context) {
	var (
		urlPrefix    string
		tag          string
		fileType     string
		saveFilePath string
		err          error
	)
	tag, _ = c.GetPostForm("type")
	fileType = c.DefaultQuery("file_type", "images")

	if fileType != "images" && fileType != "files" {
		app.Error(c, -1, fmt.Errorf("上传接口目前，仅支持图片上传和文件上传"), "")
		return
	}

	if viper.GetBool("settings.domain.getHost") {
		urlPrefix = fmt.Sprintf("http://%s/", c.Request.Host)
	} else {
		if strings.HasSuffix(viper.GetString("settings.domain.url"), "/") {
			urlPrefix = viper.GetString("settings.domain.url")
		} else {
			urlPrefix = fmt.Sprintf("http://%s/", viper.GetString("settings.domain.url"))
		}
	}

	if tag == "" {
		tag = "1"
	}

	saveFilePath = "static/uploadfile/" + fileType + "/"
	_, err = os.Stat(saveFilePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(saveFilePath, 0755)
		if err != nil {
			app.Error(c, -1, err, fmt.Sprintf("创建图片目录失败，%v", err.Error()))
			return
		}
	}

	switch tag {
	case "1": // 单图
		files, err := c.FormFile("file")
		if err != nil {
			app.Error(c, 200, errors.New(""), "图片不能为空")
			return
		}
		// 上传文件至指定目录
		guid := uuid.New().String()

		singleFile := saveFilePath + guid + utils.GetExt(files.Filename)
		_ = c.SaveUploadedFile(files, singleFile)
		app.OK(c, urlPrefix+singleFile, "上传成功")
		return
	case "2": // 多图
		files := c.Request.MultipartForm.File["file"]
		multipartFile := make([]string, len(files))
		for _, f := range files {
			guid := uuid.New().String()
			multipartFileName := saveFilePath + guid + utils.GetExt(f.Filename)
			_ = c.SaveUploadedFile(f, multipartFileName)
			multipartFile = append(multipartFile, urlPrefix+multipartFileName)
		}
		app.OK(c, multipartFile, "上传成功")
		return
	case "3": // base64
		files, _ := c.GetPostForm("file")
		ddd, _ := base64.StdEncoding.DecodeString(files)
		guid := uuid.New().String()
		_ = ioutil.WriteFile(saveFilePath+guid+".jpg", ddd, 0666)
		app.OK(c, urlPrefix+saveFilePath+guid+".jpg", "上传成功")
	default:
		app.Error(c, 200, errors.New(""), "标识不正确")
		return
	}
}
