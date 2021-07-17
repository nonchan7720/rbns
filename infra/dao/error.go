package dao

import "github.com/n-creativesystem/rbnc/infra/dao/driver/postgres"

func NewDBErr(err error) error {
	switch dialector {
	case postgreSQL:
		return postgres.NewDBErr(err)
	case mySQL:
		return err
	}
	return err
}
