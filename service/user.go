package service

import (
	"context"

	"github.com/n-creativesystem/api-rbac/domain/model"
	"github.com/n-creativesystem/api-rbac/domain/repository"
	"github.com/n-creativesystem/api-rbac/proto"
	"github.com/n-creativesystem/api-rbac/protoconv"
)

type userService struct {
	*proto.UnimplementedUserServer
	repo repository.Repository
}

var _ proto.UserServer = (*userService)(nil)

func NewUserService(repo repository.Repository) proto.UserServer {
	return &userService{repo: repo}
}

// User
func (srv *userService) Create(ctx context.Context, in *proto.UserEntity) (*proto.Empty, error) {
	orgId, err := model.NewID(in.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	roles := make([]*proto.RoleEntity, len(in.GetRoles()))
	copy(roles, in.GetRoles())
	user, err := model.NewUser(in.GetKey(), nil, nil)
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		mRole := make(model.Roles, len(roles))
		for idx, role := range roles {
			if id, err := model.NewID(role.GetId()); err == nil {
				if r, err := tx.Role().FindByID(id); err == nil {
					mRole[idx] = *r
				} else {
					return err
				}
			} else {
				return err
			}
		}
		user.AddRole(mRole...)
		userRepo := tx.User()
		_, err := userRepo.Create(orgId, user)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *userService) Delete(ctx context.Context, in *proto.UserKey) (*proto.Empty, error) {
	orgId, err := model.NewID(in.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	key, err := model.NewKey(in.GetKey())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		return tx.User().Delete(orgId, key)
	})
	return &proto.Empty{}, err
}

func (srv *userService) FindByKey(ctx context.Context, in *proto.UserKey) (*proto.UserEntity, error) {
	orgId, err := model.NewID(in.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	key, err := model.NewKey(in.GetKey())
	if err != nil {
		return nil, err
	}
	userRepo := srv.repo.NewConnection().User(ctx)
	u, err := userRepo.FindByKey(orgId, key)
	if err != nil {
		return nil, err
	}
	out := protoconv.NewUserEntityByModel(*u)
	out.OrganizationId = in.GetOrganizationId()
	return out, nil
}

func (srv *userService) AddRole(ctx context.Context, in *proto.UserRole) (*proto.Empty, error) {
	orgId, err := model.NewID(in.GetUser().GetOrganizationId())
	key, err := model.NewKey(in.GetUser().GetKey())
	roles := make([]*proto.RoleKey, len(in.GetRoles()))
	copy(roles, in.GetRoles())
	if len(roles) == 0 {
		return &proto.Empty{}, nil
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		mRole := make(model.Roles, len(roles))
		for idx, role := range roles {
			if id, err := model.NewID(role.GetId()); err == nil {
				if r, err := tx.Role().FindByID(id); err == nil {
					mRole[idx] = *r
				} else {
					return err
				}
			} else {
				return err
			}
		}
		return tx.User().AddRole(orgId, key, mRole...)
	})
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *userService) DeleteRole(ctx context.Context, in *proto.UserDeleteRole) (*proto.Empty, error) {
	orgId, err := model.NewID(in.GetUser().GetOrganizationId())
	if err != nil {
		return nil, err
	}
	key, err := model.NewKey(in.GetUser().GetKey())
	if err != nil {
		return nil, err
	}
	roleId, err := model.NewID(in.GetRole().GetId())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		return tx.User().DeleteRole(orgId, key, roleId)
	})
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}
