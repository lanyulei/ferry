package ldap

import (
	"encoding/json"
	"errors"
	"ferry/global/orm"
	"ferry/models/system"
	"ferry/pkg/logger"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

func getLdapFields() (ldapFields []map[string]string, err error) {
	var (
		settingsValue system.Settings
		contentList   []map[string]string
	)

	err = orm.Eloquent.Model(&settingsValue).Where("classify = 2").Find(&settingsValue).Error
	if err != nil {
		return
	}

	err = json.Unmarshal(settingsValue.Content, &contentList)
	if err != nil {
		return
	}

	for _, v := range contentList {
		if v["ldap_field_name"] != "" {
			ldapFields = append(ldapFields, v)
		}
	}
	return
}

func searchRequest(username string) (userInfo *ldap.Entry, err error) {
	var (
		ldapFields       []map[string]string
		cur              *ldap.SearchResult
		ldapFieldsFilter = []string{
			"dn",
		}
	)
	ldapFields, err = getLdapFields()
	for _, v := range ldapFields {
		ldapFieldsFilter = append(ldapFieldsFilter, v["ldap_field_name"])
	}
	// 用来获取查询权限的用户。如果 ldap 禁止了匿名查询，那我们就需要先用这个帐户 bind 以下才能开始查询
	if !viper.GetBool("settings.ldap.anonymousQuery") {
		err = conn.Bind(
			viper.GetString("settings.ldap.bindUserDn"),
			viper.GetString("settings.ldap.bindPwd"))
		if err != nil {
			logger.Error("用户或密码错误。", err)
			return
		}
	}

	sql := ldap.NewSearchRequest(
		viper.GetString("settings.ldap.baseDn"),
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		fmt.Sprintf("(%v=%v)", viper.GetString("settings.ldap.userField"), username),
		ldapFieldsFilter,
		nil)

	if cur, err = conn.Search(sql); err != nil {
		err = errors.New(fmt.Sprintf("在Ldap搜索用户失败, %v", err))
		logger.Error(err)
		return
	}

	if len(cur.Entries) == 0 {
		err = errors.New("未查询到对应的用户信息。")
		logger.Error(err)
		return
	}

	userInfo = cur.Entries[0]

	return
}
