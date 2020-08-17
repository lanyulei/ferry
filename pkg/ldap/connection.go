package ldap

import (
	"crypto/tls"
	"ferry/pkg/logger"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/go-ldap/ldap/v3"
)

/*
  @Author : lanyulei
*/

type Connection struct {
	Conn *ldap.Conn
}

// ldap连接
func (c *Connection) ldapConnection() (err error) {
	var ldapConn = fmt.Sprintf("%v:%v", viper.GetString("settings.ldap.host"), viper.GetString("settings.ldap.port"))

	if viper.GetInt("settings.ldap.port") == 636 {
		c.Conn, err = ldap.DialTLS(
			"tcp",
			ldapConn,
			&tls.Config{InsecureSkipVerify: true},
		)
	} else {
		c.Conn, err = ldap.Dial(
			"tcp",
			ldapConn,
		)
	}

	if err != nil {
		logger.Errorf("无法连接到ldap服务器，%v", err)
		return
	}

	//设置超时时间
	c.Conn.SetTimeout(5 * time.Second)

	return
}
