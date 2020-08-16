package ldap

import (
	"ferry/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

func (c *Connection) LdapLogin(username string, password string) (err error) {
	err = c.ldapConnection()
	if err != nil {
		return
	}
	defer c.Conn.Close()

	err = c.Conn.Bind(fmt.Sprintf("cn=%v,%v", username, viper.GetString("settings.ldap.baseDn")), password)
	if err != nil {
		logger.Error("用户或密码错误。", err)
		return
	}

	return
}
