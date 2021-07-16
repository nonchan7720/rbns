package mysql

import (
	"github.com/n-creativesystem/api-rbac/domain/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), opts...)
}

func NewDBErr(err error) error {
	return model.NewDBErr(err)
}
