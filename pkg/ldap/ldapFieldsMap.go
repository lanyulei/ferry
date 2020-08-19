package ldap

import (
	"ferry/models/system"

	"github.com/go-ldap/ldap/v3"
)

/*
  @Author : lanyulei
*/

func LdapFieldsMap(ldapUserInfo *ldap.Entry) (userInfo system.SysUser, err error) {
	var (
		ldapFields []map[string]string
	)

	ldapFields, err = getLdapFields()
	if err != nil {
		return
	}

	for _, v := range ldapFields {
		switch v["local_field_name"] {
		case "nick_name":
			userInfo.NickName = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		case "phone":
			userInfo.Phone = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		case "avatar":
			userInfo.Avatar = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		case "sex":
			userInfo.Sex = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		case "email":
			userInfo.Email = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		case "remark":
			userInfo.Remark = ldapUserInfo.GetAttributeValue(v["ldap_field_name"])
		}
	}

	return
}
