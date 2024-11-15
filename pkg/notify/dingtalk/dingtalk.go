/*
 * @Author: licon licon.ye@gmail.com
 * @LastEditors: licon licon.ye@gmail.com
 * @FilePath: /pkg/notify/dingtalk/dingtalk.go
 * @Description: 发送钉钉第三方应用消息, 需要工单系统内的用户手机号与钉钉上一致
 */
package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"ferry/pkg/logger"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// var appKey = viper.GetString("settings.dingtalk.appkey")
// var appSecret = viper.GetString("settings.dingtalk.appsecret") 
// var agentId = viper.GetInt64("settings.dingtalk.agentid") 

const (
	GetAccessTokenUrl      = "https://oapi.dingtalk.com/gettoken"
	SendWorkMsgUrl         = "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2"
	GetUserInfoByMobileUrl = "https://oapi.dingtalk.com/topapi/v2/user/getbymobile"
)

type UserDetailInfo struct {
	LoginName string `gorm:"login_name" json:"login_name"`
	UserId    string `gorm:"user_id" json:"user_id"`
}
type TokenResp struct {
	Errcode     int64  `json:"errcode"`
	AccessToken string `json:"access_token"`
	Errmsg      string `json:"errmsg"`
	ExpiresIn   int64  `json:"expires_in"`
}

type MsgBody struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}
type SendFlowMsgBody struct {
	AgentId    int64   `json:"agent_id"`
	UseridList string  `json:"userid_list"`
	Msg        MsgBody `json:"msg"`
}

type SendFlowResp struct {
	Errcode   int64  `json:"errcode"`
	TaskId    int64  `json:"task_id"`
	RequestId string `json:"request_id"`
}

type GetUserByMobileReq struct {
	Mobile string `json:"mobile"`
}

type Result struct {
	Userid string `json:"userid"`
}

type GetUserByMobileResp struct {
	Errcode   int64   `json:"errcode"`
	RequestId string  `json:"request_id"`
	Result    *Result `json:"result"`
}

func GetDingUidByMobile(token string, req GetUserByMobileReq) string {
	b, _ := json.Marshal(req)
	resp, err := http.Post(fmt.Sprintf("%s?access_token=%s", GetUserInfoByMobileUrl, token), "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Info(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var jsonBody GetUserByMobileResp
	json.Unmarshal(body, &jsonBody)
	if jsonBody.Errcode != 0 {
		logger.Info("GetDingUidByMobile error")
		return ""
	}
	if jsonBody.Result != nil {
		return jsonBody.Result.Userid
	} else {
		return ""
	}
}

func GetAccessToken(appKey, appSecret string) string {
	res, err := http.Get(fmt.Sprintf("%s?appkey=%s&appsecret=%s", GetAccessTokenUrl, appKey, appSecret))
	if err != nil {
		logger.Info(err)
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Info(err)
		return ""
	}
	var token TokenResp
	json.Unmarshal(body, &token)
	if err != nil {
		logger.Info(err)
		return ""
	}
	if token.Errcode == 0 {
		return token.AccessToken
	} else {
		logger.Info(token.Errmsg)
		return ""
	}

}

func SendFlowMsg(token string, content SendFlowMsgBody) error {
	b, _ := json.Marshal(content)
	resp, err := http.Post(fmt.Sprintf("%s?access_token=%s", SendWorkMsgUrl, token), "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Info(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var jsonBody SendFlowResp
	json.Unmarshal(body, &jsonBody)
	if jsonBody.Errcode != 0 {
		logger.Info(err)
	}
	return nil

}

func SendDingMsg(phoneList []string, url string, msgTitle string, msgCreator string, priority string, createdAt string) {
	appKey    := viper.GetString("settings.dingtalk.appkey")
	appSecret := viper.GetString("settings.dingtalk.appsecret") 
	agentId   := viper.GetInt64("settings.dingtalk.agentid")

	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("settings.redis.host"),
		Password: viper.GetString("settings.redis.pwd"),
		DB:       1, // use default DB
	})

	var token string
	tokenByte, err := client.Get("accessTokenDingtalk").Result()

	if err != nil {
		logger.Info(err)
		token = GetAccessToken(appKey, appSecret)
		err = client.Set("accessTokenDingtalk", token, 120*time.Minute).Err()
		if err != nil {
			logger.Info(err)
		}
	} else {
		token = string(tokenByte)
	}
	var userId string

	var uids []string
	var getDingUidByMobileReq GetUserByMobileReq
	for i := 0; i < len(phoneList); i++ {
		getDingUidByMobileReq.Mobile = phoneList[i]
		userId = GetDingUidByMobile(token, getDingUidByMobileReq)
		if userId != "" {
			uids = append(uids, userId)
		}
	}
	if len(uids) == 0 {
		return
	}
	uidsString := strings.Join(uids, ",")
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	var content SendFlowMsgBody
	content.AgentId = agentId
	content.UseridList = uidsString
	content.Msg.Msgtype = "markdown"
	content.Msg.Markdown.Title = "新工单提醒"
	content.Msg.Markdown.Text = "# 你有一个待办工单需要处理： " +
		"\n  **标题&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;：** " + msgTitle + "  " +
		"\n  **创建人&nbsp;&nbsp;&nbsp;&nbsp;：** " + msgCreator + "  " +
		"\n  **优先级&nbsp;&nbsp;&nbsp;&nbsp;：** " + priority + "  " +
		"\n  **创建时间：** " + createdAt + "  " +
		"\n  **当前时间：** " + nowStr + "  " +
		"\n  [去处理](" + url + ")"

	err = SendFlowMsg(token, content)
	if err != nil {
		logger.Info(err)
	}
	logger.Info("send dingtalk notify successfully")
	return
}
