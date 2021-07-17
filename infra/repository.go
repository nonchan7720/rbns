package infra

import (
	"context"

	"github.com/n-creativesystem/rbnc/domain/repository"
	"github.com/n-creativesystem/rbnc/infra/dao"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type dbRepository struct {
	driver dao.DataBase
}

// @repository
func NewRepository(driver dao.DataBase) repository.Repository {
	return &dbRepository{
		driver: driver,
	}
}

var _ repository.Repository = (*dbRepository)(nil)

func (repo *dbRepository) NewConnection() repository.Connection {
	return &dbConnection{
		driver: repo.driver,
	}
}

type dbConnection struct {
	driver dao.DataBase
}

var _ repository.Connection = (*dbConnection)(nil)

func (con *dbConnection) Permission(ctx context.Context) repository.Permission {
	return &permission{
		master: con.driver.Session(ctx),
		slave:  con.driver.SessionSlave(ctx),
	}
}

func (con *dbConnection) Role(ctx context.Context) repository.Role {
	return &role{
		master: con.driver.Session(ctx),
		slave:  con.driver.SessionSlave(ctx),
	}
}

func (con *dbConnection) Organization(ctx context.Context) repository.Organization {
	return &organization{
		master: con.driver.Session(ctx),
		slave:  con.driver.SessionSlave(ctx),
	}
}

func (con *dbConnection) User(ctx context.Context) repository.User {
	return &user{
		master: con.driver.Session(ctx),
		slave:  con.driver.SessionSlave(ctx),
	}
}

func (con *dbConnection) Transaction(ctx context.Context) repository.Tx {
	return &transaction{
		master: con.driver.Session(ctx),
		slave:  con.driver.SessionSlave(ctx),
	}
}

type transaction struct {
	master *gorm.DB
	slave  *gorm.DB
}

var _ repository.Tx = (*transaction)(nil)

func (t *transaction) Do(fn func(tx repository.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Println(err)
		}
	}()
	tx := t.master.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	tx.SkipDefaultTransaction = true
	err = fn(&dbTransaction{
		db: tx,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

type dbTransaction struct {
	db *gorm.DB
}

var _ repository.Transaction = (*dbTransaction)(nil)

func (tx *dbTransaction) Permission() repository.PermissionCommand {
	return &permission{
		master: tx.db,
		slave:  tx.db,
	}
}

func (tx *dbTransaction) Role() repository.RoleCommand {
	return &role{
		master: tx.db,
		slave:  tx.db,
	}
}

func (tx *dbTransaction) Organization() repository.OrganizationCommand {
	return &organization{
		master: tx.db,
		slave:  tx.db,
	}
}

func (tx *dbTransaction) User() repository.UserCommand {
	return &user{
		master: tx.db,
		slave:  tx.db,
	}
}
