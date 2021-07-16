package service_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/n-creativesystem/api-rbac/domain/repository"
	"github.com/n-creativesystem/api-rbac/infra"
	"github.com/n-creativesystem/api-rbac/infra/dao"
	"github.com/n-creativesystem/api-rbac/proto"
	"github.com/n-creativesystem/api-rbac/service"
	"github.com/n-creativesystem/api-rbac/tests"
	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				expecteds := []struct {
					name, description string
				}{
					{
						name:        "create:user",
						description: "create user permission",
					},
					{
						name:        "update:user",
						description: "update user permisssion",
					},
					{
						name:        "read:user",
						description: "read user permission",
					},
					{
						name:        "delete:user",
						description: "delete user permission",
					},
				}
				var repo repository.Repository = infra.NewRepository(db)
				pSrv := service.NewPermissionService(repo)
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "permissions" ("id","created_at","updated_at","name","description") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10),($11,$12,$13,$14,$15),($16,$17,$18,$19,$20)`),
				).WillReturnResult(sqlmock.NewResult(0, 4))
				mock.ExpectCommit()
				return func(t *testing.T) {
					out, err := pSrv.Create(ctx, &proto.PermissionEntities{
						Permissions: []*proto.PermissionEntity{
							{
								Name:        "create:user",
								Description: "create user permission",
							},
							{
								Name:        "update:user",
								Description: "update user permisssion",
							},
							{
								Name:        "read:user",
								Description: "read user permission",
							},
							{
								Name:        "delete:user",
								Description: "delete user permission",
							},
						},
					})
					assert.NoError(t, err)
					for idx, entity := range out.Permissions {
						expected := expecteds[idx]
						assert.NotEmpty(t, entity.GetId())
						assert.Equal(t, expected.name, entity.GetName())
						assert.Equal(t, expected.description, entity.GetDescription())
					}
				}
			},
		},
		{
			Name: "findById",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				var repo repository.Repository = infra.NewRepository(db)
				pSrv := service.NewPermissionService(repo)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" = $1`),
				).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("1", "create:user", "create user permission"))
				return func(t *testing.T) {
					out, err := pSrv.FindById(ctx, &proto.PermissionKey{Id: "1"})
					assert.NoError(t, err)
					assert.Equal(t, "create:user", out.GetName())
					assert.Equal(t, "create user permission", out.GetDescription())
				}
			},
		},
		{
			Name: "findAll",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				expecteds := []struct {
					id                string
					name, description string
				}{
					{
						id:          "1",
						name:        "create:user",
						description: "create user permission",
					},
					{
						id:          "2",
						name:        "update:user",
						description: "update user permisssion",
					},
					{
						id:          "3",
						name:        "read:user",
						description: "read user permission",
					},
					{
						id:          "4",
						name:        "delete:user",
						description: "delete user permission",
					},
				}
				row := sqlmock.NewRows([]string{"id", "name", "description"})
				for _, e := range expecteds {
					row.AddRow(e.id, e.name, e.description)
				}
				var repo repository.Repository = infra.NewRepository(db)
				pSrv := service.NewPermissionService(repo)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" ORDER BY id`),
				).WillReturnRows(row)
				return func(t *testing.T) {
					out, err := pSrv.FindAll(ctx, &proto.Empty{})
					assert.NoError(t, err)
					for idx, entity := range out.Permissions {
						expected := expecteds[idx]
						assert.Equal(t, expected.id, entity.GetId())
						assert.Equal(t, expected.name, entity.GetName())
						assert.Equal(t, expected.description, entity.GetDescription())
					}
				}
			},
		},
		{
			Name: "update",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				var repo repository.Repository = infra.NewRepository(db)
				pSrv := service.NewPermissionService(repo)
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`UPDATE "permissions" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "permissions"."id" = $4`),
				).WithArgs(sqlmock.AnyArg(), "read:permission", "read user permission", "1").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				return func(t *testing.T) {
					_, err := pSrv.Update(ctx, &proto.PermissionEntity{
						Id:          "1",
						Name:        "read:permission",
						Description: "read user permission",
					})
					assert.NoError(t, err)
				}
			},
		},
	}

	cases.Run(t)
}
