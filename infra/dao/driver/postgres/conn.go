package postgres

import (
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/n-creativesystem/rbnc/domain/model"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), opts...)
}

func NewDBErr(err error) error {
	if err == nil {
		return nil
	}
	if pgError, ok := err.(*pgconn.PgError); ok {
		switch pgError.SQLState() {
		case "23505":
			return status.Error(http.StatusConflict, model.ErrAlreadyExists.Error())
		}
	}
	if err == gorm.ErrRecordNotFound {
		return status.Error(http.StatusNotFound, model.ErrNoData.Error())
	}
	return model.NewDBErr(err)

}
