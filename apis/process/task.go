package process

//import (
//	"ferry/pkg/pagination"
//	"ferry/pkg/response/code"
//	. "ferry/pkg/response/response"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"strings"
//
//	"github.com/gin-gonic/gin"
//	uuid "github.com/satori/go.uuid"
//	"github.com/spf13/viper"
//)
//
///*
//  @Author : lanyulei
//*/
//
//// 任务列表
//func TaskList(c *gin.Context) {
//	var (
//		err         error
//		pageValue   pagination.ListRequest
//		taskName    string
//		taskData    []map[string]interface{}
//		total_count int
//	)
//	taskName = c.DefaultQuery("name", "")
//
//	err = c.ShouldBind(&pageValue)
//	if err != nil {
//		Response(c, code.BindError, nil, err.Error())
//		return
//	}
//
//	getFileDetails := func(fn string) map[string]interface{} {
//		file := make(map[string]interface{})
//		fileClassify := strings.Split(fn, ".")
//		fileDetails := strings.Split(fileClassify[0], "-")
//		switch fileClassify[1] {
//		case "py":
//			file["classify"] = "Python"
//		case "sh":
//			file["classify"] = "Shell"
//		default:
//			file["classify"] = "Other"
//		}
//		if len(fileDetails) == 3 {
//			file["name"] = fileDetails[0]
//			file["uuid"] = fileDetails[1]
//			file["creator"] = fileDetails[2]
//		}
//		file["full_name"] = fn
//		return file
//	}
//	files, _ := ioutil.ReadDir(viper.GetString("script.path"))
//	var endIndex int
//	if taskName != "" {
//		for _, f := range files {
//			if strings.Contains(strings.Split(f.Name(), "-")[0], taskName) {
//				taskData = append(taskData, getFileDetails(f.Name()))
//			}
//		}
//		total_count = len(taskData)
//		if pageValue.Page*pageValue.PerPage > len(taskData) {
//			endIndex = len(taskData)
//		} else {
//			endIndex = pageValue.Page * pageValue.PerPage
//		}
//		taskData = taskData[(pageValue.Page-1)*pageValue.PerPage : endIndex]
//	} else {
//		if pageValue.Page*pageValue.PerPage > len(files) {
//			endIndex = len(files)
//		} else {
//			endIndex = pageValue.Page * pageValue.PerPage
//		}
//		for _, f := range files[(pageValue.Page-1)*pageValue.PerPage : endIndex] {
//			taskData = append(taskData, getFileDetails(f.Name()))
//		}
//		total_count = len(files)
//	}
//
//	Response(c, nil, map[string]interface{}{
//		"data":        taskData,
//		"page":        pageValue.Page,
//		"per_page":    pageValue.PerPage,
//		"total_count": total_count,
//	}, "")
//}
//
//// 创建任务
//func CreateTask(c *gin.Context) {
//	type Task struct {
//		Name     string `json:"name"`
//		Classify string `json:"classify"`
//		Content  string `json:"content"`
//	}
//
//	var (
//		err       error
//		taskValue Task
//	)
//
//	err = c.ShouldBind(&taskValue)
//	if err != nil {
//		Response(c, code.BindError, nil, err.Error())
//		return
//	}
//
//	uuidValue := uuid.Must(uuid.NewV4(), err)
//	fileName := fmt.Sprintf("%v/%v-%v-%v",
//		viper.GetString("script.path"),
//		taskValue.Name,
//		strings.Split(uuidValue.String(), "-")[4],
//		c.GetString("username"),
//	)
//	if taskValue.Classify == "python" {
//		fileName = fileName + ".py"
//	} else if taskValue.Classify == "shell" {
//		fileName = fileName + ".sh"
//	}
//
//	err = ioutil.WriteFile(fileName, []byte(taskValue.Content), 0666)
//	if err != nil {
//		Response(c, code.BindError, nil, fmt.Sprintf("创建任务脚本失败: %v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 更新任务
//func UpdateTask(c *gin.Context) {
//	type fileDetails struct {
//		Name     string `json:"name"`
//		FullName string `json:"full_name"`
//		Classify string `json:"classify"`
//		Content  string `json:"content"`
//	}
//
//	var (
//		err  error
//		file fileDetails
//	)
//
//	err = c.ShouldBind(&file)
//	if err != nil {
//		Response(c, code.BindError, nil, "")
//		return
//	}
//
//	fullNameList := strings.Split(file.FullName, "-")
//	if fullNameList[0] != file.Name {
//		fullNameList[0] = file.Name
//	}
//	var suffixName string
//	if strings.ToLower(file.Classify) == "python" {
//		suffixName = ".py"
//	} else if strings.ToLower(file.Classify) == "shell" {
//		suffixName = ".sh"
//	}
//
//	if fullNameList[len(fullNameList)-1][len(fullNameList[len(fullNameList)-1])-3:len(fullNameList[len(fullNameList)-1])] != suffixName {
//		tList := strings.Split(fullNameList[len(fullNameList)-1], ".")
//		tList[len(tList)-1] = suffixName[1:len(suffixName)]
//		fullNameList[len(fullNameList)-1] = strings.Join(tList, ".")
//	}
//
//	fileFullName := strings.Join(fullNameList, "-")
//
//	// 修改文件内容
//	err = ioutil.WriteFile(fmt.Sprintf("%v/%v", viper.GetString("script.path"), fileFullName), []byte(file.Content), 0666)
//	if err != nil {
//		Response(c, code.BindError, nil, fmt.Sprintf("更新脚本文件失败，%v", err.Error()))
//		return
//	}
//
//	// 修改文件名称
//	err = os.Rename(
//		fmt.Sprintf("%v/%v", viper.GetString("script.path"), file.FullName),
//		fmt.Sprintf("%v/%v", viper.GetString("script.path"), fileFullName),
//	)
//	if err != nil {
//		Response(c, code.BindError, nil, fmt.Sprintf("更改脚本文件名称失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 删除任务
//func DeleteTask(c *gin.Context) {
//	fullName := c.DefaultQuery("full_name", "")
//	if fullName == "" {
//		Response(c, code.InternalServerError, nil, "参数不正确，请确定参数full_name是否传递")
//		return
//	}
//
//	err := os.Remove(fmt.Sprintf("%v/%v", viper.GetString("script.path"), fullName))
//	if err != nil {
//		Response(c, code.DeleteError, nil, fmt.Sprintf("删除文件失败，%v", err.Error()))
//		return
//	}
//
//	Response(c, nil, nil, "")
//}
//
//// 任务详情
//func TaskDetails(c *gin.Context) {
//	var (
//		err      error
//		fileName string
//		content  []byte
//	)
//
//	fileName = c.DefaultQuery("file_name", "")
//	if fileName == "" {
//		Response(c, code.ParamError, nil, "参数不正确，请确认file_name参数是否存在")
//		return
//	}
//
//	content, err = ioutil.ReadFile(fmt.Sprintf("%v/%v", viper.GetString("script.path"), fileName))
//	if err != nil {
//		return
//	}
//
//	Response(c, nil, string(content), "")
//}
