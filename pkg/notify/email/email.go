package email

/*
  @Author : lanyulei
  @Desc : 发送邮件
*/

import (
	"ferry/pkg/logger"
	"strconv"

	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

func server(mailTo []string, ccTo []string, subject, body string, args ...string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": viper.GetString("settings.email.user"),
		"pass": viper.GetString("settings.email.pass"),
		"host": viper.GetString("settings.email.host"),
		"port": viper.GetString("settings.email.port"),
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], viper.GetString("settings.email.alias"))) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", mailTo...)                                                                    //发送给多个用户
	m.SetHeader("Cc", ccTo...)                                                                      //发送给多个用户
	m.SetHeader("Subject", subject)                                                                 //设置邮件主题
	m.SetBody("text/html", body)                                                                    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err

}

func SendMail(mailTo []string, ccTo []string, subject, body string) {
	err := server(mailTo, ccTo, subject, body)
	if err != nil {
		logger.Info(err)
		return
	}
	logger.Info("send successfully")
}
