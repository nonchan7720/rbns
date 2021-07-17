package service

import (
	"context"

	"github.com/n-creativesystem/rbnc/domain/model"
	"github.com/n-creativesystem/rbnc/domain/repository"
	"github.com/n-creativesystem/rbnc/proto"
	"github.com/n-creativesystem/rbnc/protoconv"
)

type permissionService struct {
	*proto.UnimplementedPermissionServer
	repo repository.Repository
}

var _ proto.PermissionServer = (*permissionService)(nil)

func NewPermissionService(repo repository.Repository) proto.PermissionServer {
	return &permissionService{repo: repo}
}

// Permission
func (srv *permissionService) Create(ctx context.Context, in *proto.PermissionEntities) (*proto.PermissionEntities, error) {
	var out *proto.PermissionEntities
	inPermissions := make([]*proto.PermissionEntity, len(in.GetPermissions()))
	copy(inPermissions, in.GetPermissions())
	if len(inPermissions) == 0 {
		return &proto.PermissionEntities{
			Permissions: make([]*proto.PermissionEntity, 0),
		}, nil
	}
	names := make([]model.Name, len(inPermissions))
	descriptions := make([]string, len(inPermissions))
	for idx, permission := range inPermissions {
		var err error
		names[idx], err = model.NewName(permission.Name)
		if err != nil {
			return nil, err
		}
		descriptions[idx] = permission.Description
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err := tx.Do(func(tx repository.Transaction) error {
		permissionRepo := tx.Permission()
		permissions, err := permissionRepo.CreateBatch(names, descriptions)
		if err != nil {
			return err
		}
		out = &proto.PermissionEntities{
			Permissions: make([]*proto.PermissionEntity, len(permissions)),
		}
		for idx, permission := range permissions {
			out.Permissions[idx] = protoconv.NewPermissionEntityByModel(*permission)
		}
		return nil
	})
	return out, err
}

func (srv *permissionService) FindById(ctx context.Context, in *proto.PermissionKey) (*proto.PermissionEntity, error) {
	permissionRepo := srv.repo.NewConnection().Permission(ctx)
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	permission, err := permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return protoconv.NewPermissionEntityByModel(*permission), nil
}

func (srv *permissionService) FindAll(ctx context.Context, in *proto.Empty) (*proto.PermissionEntities, error) {
	permissionRepo := srv.repo.NewConnection().Permission(ctx)
	permissions, err := permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	out := &proto.PermissionEntities{
		Permissions: make([]*proto.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		out.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return out, nil
}

func (srv *permissionService) Update(ctx context.Context, in *proto.PermissionEntity) (*proto.Empty, error) {
	p, err := model.NewPermission(in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Permission().Update(p) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, err
}

func (srv *permissionService) Delete(ctx context.Context, in *proto.PermissionKey) (*proto.Empty, error) {
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Permission().Delete(id) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *permissionService) Check(ctx context.Context, in *proto.PermissionCheckRequest) (*proto.PermissionCheckResult, error) {
	result := &proto.PermissionCheckResult{
		Result:  false,
		Message: "",
	}
	con := srv.repo.NewConnection()
	permissionName, err := model.NewName(in.GetPrermissionName())
	if err != nil {
		result.Message = err.Error()
		return result, err
	}
	organizationName, err := model.NewName(in.GetOrganizationName())
	if err != nil {
		result.Message = err.Error()
		return result, err
	}
	userKey, err := model.NewKey(in.GetUserKey())
	if err != nil {
		return result, err
	}
	org, err := con.Organization(ctx).FindByName(organizationName)
	if err != nil {
		result.Message = err.Error()
		return result, err
	}
	if u, ok := org.IsContainsUsers(userKey); !ok {
		result.Message = model.ErrNoData.Error()
		return result, model.ErrNoData
	} else {
		if u.IsContainsPermission(permissionName) {
			result.Result = true
			return result, nil
		}
	}
	result.Result = false
	result.Message = model.ErrNoData.Error()
	return result, model.ErrNoData
}
