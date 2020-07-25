package notify

import (
	"ferry/models/system"
	"ferry/pkg/notify/email"
)

/*
  @Author : lanyulei
  @同时发送多种通知方式
*/

func SendNotify(classify []int, sendTo interface{}, subject, body string) {
	var (
		emailList []string
	)
	for _, c := range classify {
		switch c {
		case 1: // 邮件
			for _, user := range sendTo.(map[string]interface{})["userList"].([]system.SysUser) {
				emailList = append(emailList, user.Email)
			}
			go email.SendMail(emailList, subject, body)
		}
	}
}
