package mycasbin

import (
	"ferry/global/orm"
	"ferry/pkg/logger"
	"ferry/tools/config"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-sql-driver/mysql"
)

var _ endpoint.Middleware

func Casbin() (*casbin.Enforcer, error) {
	conn := orm.MysqlConn
	Apter, err := gormadapter.NewAdapter(config.DatabaseConfig.Dbtype, conn, true)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("config/rbac_model.conf", Apter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		logger.Infof("casbin rbac_model or policy init error, message: %v \r\n", err.Error())
		return nil, err
	}
}
