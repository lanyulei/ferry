package system

import (
	"errors"
	"ferry/global/orm"
	"ferry/pkg/logger"
	"ferry/tools"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

/*
  @Author : lanyulei
*/

// User
type User struct {
	// key
	IdentityKey string
	// 用户名
	UserName  string
	FirstName string
	LastName  string
	// 角色
	Role string
}

type UserName struct {
	Username string `gorm:"type:varchar(64)" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"type:varchar(128)" json:"password"`
}

type LoginM struct {
	UserName
	PassWord
}

type SysUserId struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

type SysUserB struct {
	NickName string `gorm:"type:varchar(128)" json:"nickName"` // 昵称
	Phone    string `gorm:"type:varchar(11)" json:"phone"`     // 手机号
	RoleId   int    `gorm:"type:int(11)" json:"roleId"`        // 角色编码
	Salt     string `gorm:"type:varchar(255)" json:"salt"`     //盐
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`   //头像
	Sex      string `gorm:"type:varchar(255)" json:"sex"`      //性别
	Email    string `gorm:"type:varchar(128)" json:"email"`    //邮箱
	DeptId   int    `gorm:"type:int(11)" json:"deptId"`        //部门编码
	PostId   int    `gorm:"type:int(11)" json:"postId"`        //职位编码
	CreateBy string `gorm:"type:varchar(128)" json:"createBy"` //
	UpdateBy string `gorm:"type:varchar(128)" json:"updateBy"` //
	Remark   string `gorm:"type:varchar(255)" json:"remark"`   //备注
	Status   string `gorm:"type:int(1);" json:"status"`
	Params   string `gorm:"-" json:"params"`
	BaseModel
}

type SysUser struct {
	SysUserId
	SysUserB
	LoginM
}

func (SysUser) TableName() string {
	return "sys_user"
}

type SysUserPwd struct {
	OldPassword  string `json:"oldPassword" form:"oldPassword"`
	NewPassword  string `json:"newPassword" form:"newPassword"`
	PasswordType int    `json:"passwordType" form:"passwordType"`
}

type SysUserPage struct {
	SysUserId
	SysUserB
	LoginM
	DeptName string `gorm:"-" json:"deptName"`
}

type SysUserView struct {
	SysUserId
	SysUserB
	LoginM
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

// 获取用户数据
func (e *SysUser) Get() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}

	SysUserView.Password = ""
	return
}

func (e *SysUser) GetUserInfo() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetList() (SysUserView []SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.Find(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetPage(pageSize int, pageIndex int) ([]SysUserPage, int, error) {
	var (
		doc   []SysUserPage
		count int
	)
	table := orm.Eloquent.Select("sys_user.*,sys_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join sys_dept on sys_dept.dept_id = sys_user.dept_id")

	if e.Username != "" {
		table = table.Where("sys_user.username like ?", "%"+e.Username+"%")
	}

	if e.NickName != "" {
		table = table.Where("sys_user.nick_name like ?", "%"+e.NickName+"%")
	}

	if e.Status != "" {
		table = table.Where("sys_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("sys_user.phone like ?", "%"+e.Phone+"%")
	}

	if e.DeptId != 0 {
		table = table.Where("sys_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+tools.IntToString(e.DeptId)+"%")
	}

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("sys_user.delete_time IS NULL").Count(&count)
	return doc, count, nil
}

//加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

//添加
func (e SysUser) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table(e.TableName()).Where("username = ? and `delete_time` IS NULL", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	return
}

//修改
func (e *SysUser) Update(id int) (update SysUser, err error) {
	if e.Password != "" {
		if err = e.Encrypt(); err != nil {
			return
		}
	}
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if e.RoleId == 0 {
		e.RoleId = update.RoleId
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("user_id in (?)", id).Delete(&SysUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysUser) SetPwd(pwd SysUserPwd) (Result bool, err error) {
	user, err := e.GetUserInfo()
	if err != nil {
		tools.HasError(err, "获取用户数据失败(代码202)", 500)
	}
	_, err = tools.CompareHashAndPassword(user.Password, pwd.OldPassword)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			tools.HasError(err, "密码错误(代码202)", 500)
		}
		logger.Info(err)
		return
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	tools.HasError(err, "更新密码失败(代码202)", 500)
	return
}
