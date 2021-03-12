package database

import (
	"bytes"
	"ferry/global/orm"
	"ferry/pkg/logger"
	"ferry/tools/config"
	"strconv"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
)

var (
	DbType   string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
)

func (e *Mysql) Setup() {

	var err error
	var db Database

	db = new(Mysql)
	orm.MysqlConn = db.GetConnect()
	orm.Eloquent, err = db.Open(DbType, orm.MysqlConn)

	if err != nil {
		logger.Fatalf("%s connect error %v", DbType, err)
	} else {
		logger.Infof("%s connect success!", DbType)
	}

	if orm.Eloquent.Error != nil {
		logger.Fatalf("database error %v", orm.Eloquent.Error)
	}

	// 是否开启详细日志记录
	orm.Eloquent.LogMode(viper.GetBool("settings.gorm.logMode"))

	// 设置最大打开连接数
	orm.Eloquent.DB().SetMaxOpenConns(viper.GetInt("settings.gorm.maxOpenConn"))

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	orm.Eloquent.DB().SetMaxIdleConns(viper.GetInt("settings.gorm.maxIdleConn"))
}

type Mysql struct {
}

func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

func (e *Mysql) GetConnect() string {

	DbType = config.DatabaseConfig.Dbtype
	Host = config.DatabaseConfig.Host
	Port = config.DatabaseConfig.Port
	Name = config.DatabaseConfig.Name
	Username = config.DatabaseConfig.Username
	Password = config.DatabaseConfig.Password

	var conn bytes.Buffer
	conn.WriteString(Username)
	conn.WriteString(":")
	conn.WriteString(Password)
	conn.WriteString("@tcp(")
	conn.WriteString(Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(Port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(Name)
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=10000ms")
	return conn.String()
}
