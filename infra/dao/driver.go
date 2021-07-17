package dao

import (
	"errors"

	"github.com/n-creativesystem/rbnc/infra/dao/driver/mysql"
	"github.com/n-creativesystem/rbnc/infra/dao/driver/postgres"
	"gorm.io/gorm"
)

func newDriver(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	switch dialector {
	case postgreSQL:
		return postgres.New(dsn, opts...)
	case mySQL:
		return mysql.New(dsn, opts...)
	default:
		return nil, errors.New("can not find this dialector")
	}
}
