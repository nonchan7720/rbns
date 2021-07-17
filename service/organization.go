package service

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/protoconv"
)

type organizationService struct {
	*proto.UnimplementedOrganizationServer
	repo repository.Repository
}

var _ proto.OrganizationServer = (*organizationService)(nil)

func NewOrganizationService(repo repository.Repository) proto.OrganizationServer {
	return &organizationService{repo: repo}
}

// Organization
func (srv *organizationService) Create(ctx context.Context, in *proto.OrganizationEntity) (*proto.OrganizationEntity, error) {
	var out *proto.OrganizationEntity
	orgName, err := model.NewName(in.GetName())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		orgRepo := tx.Organization()
		org, err := orgRepo.Create(orgName, in.Description)
		if err != nil {
			return err
		}
		out = &proto.OrganizationEntity{
			Id:          *org.GetID(),
			Name:        *org.GetName(),
			Description: org.GetDescription(),
		}
		return nil
	})
	return out, err
}

func (srv *organizationService) FindById(ctx context.Context, in *proto.OrganizationKey) (*proto.OrganizationEntity, error) {
	orgRepo := srv.repo.NewConnection().Organization(ctx)
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	organization, err := orgRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return protoconv.NewOrganizationEntityByModel(*organization), nil
}

func (srv *organizationService) FindAll(ctx context.Context, in *proto.Empty) (*proto.OrganizationEntities, error) {
	orgRepo := srv.repo.NewConnection().Organization(ctx)
	organizations, err := orgRepo.FindAll()
	if err != nil {
		return nil, err
	}
	protoOrganizations := make([]*proto.OrganizationEntity, len(organizations))
	for idx, organization := range organizations {
		protoOrganizations[idx] = protoconv.NewOrganizationEntityByModel(organization)
	}
	out := &proto.OrganizationEntities{
		Organizations: protoOrganizations,
	}
	return out, nil
}

func (srv *organizationService) Update(ctx context.Context, in *proto.OrganizationUpdateEntity) (*proto.Empty, error) {
	mOrg, err := model.NewOrganization(in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Organization().Update(mOrg) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (srv *organizationService) Delete(ctx context.Context, in *proto.OrganizationKey) (*proto.Empty, error) {
	id, err := model.NewID(in.GetId())
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Organization().Delete(id) }); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}
