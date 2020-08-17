package ldap

import (
	"ferry/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

func LdapLogin(username string, password string) (err error) {
	err = ldapConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	err = conn.Bind(fmt.Sprintf("cn=%v,%v", username, viper.GetString("settings.ldap.baseDn")), password)
	if err != nil {
		logger.Error("用户或密码错误。", err)
		return
	}

	return
}
