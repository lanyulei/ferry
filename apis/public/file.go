package public

import (
	"bytes"
	"encoding/base64"
	"errors"
	"ferry/pkg/logger"
	"ferry/pkg/utils"
	"ferry/tools/app"
	"ferry/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
		urlPrefix string
		tag       string
	)
	tag, _ = c.GetPostForm("type")

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
		app.Error(c, 200, errors.New(""), "缺少标识")
		return
	} else {
		switch tag {
		case "1": // 单图
			files, err := c.FormFile("file")
			if err != nil {
				app.Error(c, 200, errors.New(""), "图片不能为空")
				return
			}
			// 上传文件至指定目录
			name, err := saveUploadFile(files)
			if err != nil {
				app.Error(c, 200, err, "上传失败")
				return
			}
			app.OK(c, urlPrefix+"upload/"+name, "上传成功")
			return
		case "2": // 多图
			form, err := c.MultipartForm()
			if err != nil {
				app.Error(c, 200, err, "上传失败")
				return
			}
			files := form.File["file"]
			multipartFile := make([]string, 0)
			for _, f := range files {
				name, err := saveUploadFile(f)
				if err != nil {
					app.Error(c, 200, err, "上传失败")
					return
				}
				multipartFileName := urlPrefix + "upload/" + name
				multipartFile = append(multipartFile, multipartFileName)
			}
			app.OK(c, multipartFile, "上传成功")
			return
		case "3": // base64
			//TODO: 应当根据header(data:image/png;base64)进行保存, 不是全部都保存为jpg
			files, _ := c.GetPostForm("file")
			ddd, _ := base64.StdEncoding.DecodeString(files)
			r := bytes.NewReader(ddd)
			name, err := saveFile(r, ".jpg")
			if err != nil {
				app.Error(c, 200, err, "上传失败")
				return
			}
			app.OK(c, urlPrefix+"upload/"+name, "上传成功")
		}
	}
}

func saveUploadFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return saveFile(f, utils.GetExt(file.Filename))
}

//saveFile 返回name为保存后的文件名称，如果err不为nil, 则name为空
func saveFile(r io.Reader, ext string) (name string, err error) {
	var (
		f *os.File
	)

	//检查上传目录是否存在,不存在则创建
	dir := config.ApplicationConfig.Upload
	if !isExistDir(dir) {
		if err = os.MkdirAll(dir, 0600); err != nil {
			logger.Debug("dir", dir)
			return "", err
		}
	}

	//生成一个文件名称
	guid := uuid.New().String()
	name = fmt.Sprintf("%s%s", guid, ext)

	//保存文件内容
	path := filepath.Join(dir, name)
	f, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = io.Copy(f, r); err != nil {
		return "", err
	}
	return name, nil
}

//isExistDir 检查给定的目录是否存在
func isExistDir(s string) bool {
	info, err := os.Stat(s)
	return err == nil && info.IsDir()
}
