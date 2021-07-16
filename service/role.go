package service

import (
	"context"

	"github.com/n-creativesystem/api-rbac/domain/model"
	"github.com/n-creativesystem/api-rbac/domain/repository"
	"github.com/n-creativesystem/api-rbac/proto"
	"github.com/n-creativesystem/api-rbac/protoconv"
)

type roleService struct {
	*proto.UnimplementedRoleServer
	repo repository.Repository
}

var _ proto.RoleServer = (*roleService)(nil)

func NewRoleService(repo repository.Repository) proto.RoleServer {
	return &roleService{repo: repo}
}

func (srv *roleService) Create(ctx context.Context, in *proto.RoleEntities) (*proto.RoleEntities, error) {
	var out *proto.RoleEntities
	inRoles := make([]*proto.RoleEntity, len(in.GetRoles()))
	copy(inRoles, in.GetRoles())
	if len(inRoles) == 0 {
		return &proto.RoleEntities{
			Roles: make([]*proto.RoleEntity, 0),
		}, nil
	}
	names := make([]model.Name, len(inRoles))
	descriptions := make([]string, len(inRoles))
	for idx, role := range inRoles {
		var err error
		names[idx], err = model.NewName(role.Name)
		if err != nil {
			return nil, err
		}
		descriptions[idx] = role.Description
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err := tx.Do(func(tx repository.Transaction) error {
		roleRepo := tx.Role()
		roles, err := roleRepo.CreateBatch(names, descriptions)
		if err != nil {
			return err
		}
		out = &proto.RoleEntities{
			Roles: make([]*proto.RoleEntity, len(roles)),
		}
		for idx, role := range roles {
			out.Roles[idx] = protoconv.NewRoleEntityByModel(*role)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (srv *roleService) FindById(ctx context.Context, in *proto.RoleKey) (*proto.RoleEntity, error) {
	roleId, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	role, err := srv.repo.NewConnection().Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	return protoconv.NewRoleEntityByModel(*role), nil
}

func (srv *roleService) FindAll(ctx context.Context, in *proto.Empty) (*proto.RoleEntities, error) {
	roleRepo := srv.repo.NewConnection().Role(ctx)
	roles, err := roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	entities := &proto.RoleEntities{
		Roles: make([]*proto.RoleEntity, len(roles)),
	}
	for i, role := range roles {
		entities.Roles[i] = protoconv.NewRoleEntityByModel(role)
	}
	return entities, nil
}

func (srv *roleService) Update(ctx context.Context, in *proto.RoleUpdateEntity) (*proto.Empty, error) {
	mRole, err := model.NewRole(in.GetId(), in.GetName(), in.GetDescription(), nil)
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Role().Update(mRole) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *roleService) Delete(ctx context.Context, in *proto.RoleKey) (*proto.Empty, error) {
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Role().Delete(id) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *roleService) GetPermissions(ctx context.Context, in *proto.RoleKey) (*proto.PermissionEntities, error) {
	roleId, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	role, err := srv.repo.NewConnection().Role(ctx).FindByID(roleId)
	if err != nil {
		return nil, err
	}
	permissions := role.GetPermissions().Copy()
	res := proto.PermissionEntities{
		Permissions: make([]*proto.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		res.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return &res, nil
}

func (srv *roleService) AddPermissions(ctx context.Context, in *proto.RoleReleationPermissions) (*proto.Empty, error) {
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	inPermissions := make([]*proto.PermissionKey, len(in.GetPermissions()))
	copy(inPermissions, in.GetPermissions())
	if len(inPermissions) == 0 {
		return &proto.Empty{}, nil
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		permissions := model.Permissions{}
		for _, permission := range inPermissions {
			if pId, err := model.NewID(permission.GetId()); err == nil {
				if p, err := tx.Permission().FindByID(pId); err != nil {
					return err
				} else {
					permissions = append(permissions, *p)
				}
			} else {
				return err
			}
		}
		return tx.Role().AddPermission(id, permissions)
	})
	return &proto.Empty{}, err
}

func (srv *roleService) DeletePermission(ctx context.Context, in *proto.RoleReleationPermission) (*proto.Empty, error) {
	roleId, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	permissionId, err := model.NewID(in.GetPermissionId())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Role().DeletePermission(roleId, permissionId) }); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
