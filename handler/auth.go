package handler

import (
	"errors"
	"ferry/global/orm"
	"ferry/models/system"
	jwt "ferry/pkg/jwtauth"
	ldap1 "ferry/pkg/ldap"
	"ferry/pkg/logger"
	"ferry/tools"
	"fmt"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/mssola/user_agent"
)

var store = base64Captcha.DefaultMemStore

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(system.SysUser)
		r, _ := v["role"].(system.SysRole)
		return jwt.MapClaims{
			jwt.IdentityKey: u.UserId,
			jwt.RoleIdKey:   r.RoleId,
			jwt.RoleKey:     r.RoleKey,
			jwt.NiceKey:     u.Username,
			jwt.RoleNameKey: r.RoleName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["nice"],
		"RoleKey":     claims["rolekey"],
		"UserId":      claims["identity"],
		"RoleIds":     claims["roleid"],
	}
}

// @Summary 登陆
// @Description 获取token
// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Param username body models.Login  true "Add account"
// @Success 200 {string} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
// @Router /login [post]
func Authenticator(c *gin.Context) (interface{}, error) {
	var (
		err           error
		loginVal      system.Login
		loginLog      system.LoginLog
		roleValue     system.SysRole
		authUserCount int
		addUserInfo   system.SysUser
		ldapUserInfo  *ldap.Entry
	)

	ua := user_agent.New(c.Request.UserAgent())
	loginLog.Ipaddr = c.ClientIP()
	location := tools.GetLocation(c.ClientIP())
	loginLog.LoginLocation = location
	loginLog.LoginTime = tools.GetCurrntTime()
	loginLog.Status = "0"
	loginLog.Remark = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	loginLog.Browser = browserName + " " + browserVersion
	loginLog.Os = ua.OS()
	loginLog.Msg = "登录成功"
	loginLog.Platform = ua.Platform()

	// 获取前端过来的数据
	if err := c.ShouldBind(&loginVal); err != nil {
		loginLog.Status = "1"
		loginLog.Msg = "数据解析失败"
		loginLog.Username = loginVal.Username
		_, _ = loginLog.Create()
		return nil, jwt.ErrMissingLoginValues
	}
	loginLog.Username = loginVal.Username

	// 校验验证码
	if !store.Verify(loginVal.UUID, loginVal.Code, true) {
		loginLog.Status = "1"
		loginLog.Msg = "验证码错误"
		_, _ = loginLog.Create()
		return nil, jwt.ErrInvalidVerificationode
	}

	// ldap 验证
	if loginVal.LoginType == 1 {
		// ldap登陆
		ldapUserInfo, err = ldap1.LdapLogin(loginVal.Username, loginVal.Password)
		if err != nil {
			return nil, err
		}
		// 2. 将ldap用户信息写入到用户数据表中
		err = orm.Eloquent.Model(&system.SysUser{}).
			Where("username = ?", loginVal.Username).
			Count(&authUserCount).Error
		if err != nil {
			return nil, errors.New(fmt.Sprintf("查询用户失败，%v", err))
		}
		addUserInfo, err = ldap1.LdapFieldsMap(ldapUserInfo)
		if err != nil {
			return nil, fmt.Errorf("ldap映射本地字段失败，%v", err.Error())
		}
		if authUserCount == 0 {
			addUserInfo.Username = loginVal.Username
			// 获取默认权限ID
			err = orm.Eloquent.Model(&system.SysRole{}).Where("role_key = 'common'").Find(&roleValue).Error
			if err != nil {
				return nil, errors.New(fmt.Sprintf("查询角色失败，%v", err))
			}
			addUserInfo.RoleId = roleValue.RoleId // 绑定通用角色
			addUserInfo.Status = "0"
			addUserInfo.CreatedAt = time.Now()
			addUserInfo.UpdatedAt = time.Now()
			if addUserInfo.Sex == "" {
				addUserInfo.Sex = "0"
			}
			err = orm.Eloquent.Create(&addUserInfo).Error
			if err != nil {
				return nil, errors.New(fmt.Sprintf("创建本地用户失败，%v", err))
			}
		}
	}

	user, role, e := loginVal.GetUser()
	if e == nil {
		_, _ = loginLog.Create()

		if user.Status == "1" {
			return nil, errors.New("用户已被禁用。")
		}

		return map[string]interface{}{"user": user, "role": role}, nil
	} else {
		loginLog.Status = "1"
		loginLog.Msg = "登录失败"
		_, _ = loginLog.Create()
		logger.Info(e.Error())
	}

	return nil, jwt.ErrFailedAuthentication
}

// @Summary 退出登录
// @Description 获取token
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /logout [post]
// @Security
func LogOut(c *gin.Context) {
	var loginlog system.LoginLog
	ua := user_agent.New(c.Request.UserAgent())
	loginlog.Ipaddr = c.ClientIP()
	location := tools.GetLocation(c.ClientIP())
	loginlog.LoginLocation = location
	loginlog.LoginTime = tools.GetCurrntTime()
	loginlog.Status = "0"
	loginlog.Remark = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	loginlog.Browser = browserName + " " + browserVersion
	loginlog.Os = ua.OS()
	loginlog.Platform = ua.Platform()
	loginlog.Username = tools.GetUserName(c)
	loginlog.Msg = "退出成功"
	_, _ = loginlog.Create()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})

}

func Authorizator(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(system.SysUser)
		r, _ := v["role"].(system.SysRole)
		c.Set("role", r.RoleName)
		c.Set("roleIds", r.RoleId)
		c.Set("userId", u.UserId)
		c.Set("userName", u.UserName)

		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
