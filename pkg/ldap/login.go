package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

/*
  @Author : lanyulei
*/

func LdapLogin(username string, password string) (userInfo *ldap.Entry, err error) {
	err = ldapConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	userInfo, err = searchRequest(username)
	if err != nil {
		return
	}

	err = conn.Bind(userInfo.DN, password)
	if err != nil {
		return nil, fmt.Errorf("用户或密码不正确。")
	}

	return
}
