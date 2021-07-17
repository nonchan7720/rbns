package infra

import (
	"context"

	"github.com/n-creativesystem/rbnc/domain/repository"
	"github.com/n-creativesystem/rbnc/infra/dao"
	"github.com/n-creativesystem/rbnc/infra/entity"
	"gorm.io/gorm"
)

type apiKey struct {
	master *gorm.DB
	slave  *gorm.DB
}

var _ repository.ApiKey = (*apiKey)(nil)

func NewAuth(db dao.DataBase) repository.ApiKey {
	return &apiKey{
		master: db.Session(context.Background()),
		slave:  db.SessionSlave(context.Background()),
	}
}

func (r *apiKey) Generate() (string, error) {
	t := entity.ApiKey{}
	t.Generate()
	err := r.master.Create(&t).Error
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

func (r *apiKey) Get() string {
	var t entity.ApiKey
	r.slave.Find(&t)
	return t.AccessToken
}
