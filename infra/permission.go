package infra

import (
	"github.com/n-creativesystem/rbnc/domain/model"
	"github.com/n-creativesystem/rbnc/domain/repository"
	"github.com/n-creativesystem/rbnc/infra/entity"
	"gorm.io/gorm"
)

type permission struct {
	master *gorm.DB
	slave  *gorm.DB
}

var (
	_ repository.Permission        = (*permission)(nil)
	_ repository.PermissionCommand = (*permission)(nil)
)

func (r *permission) FindAll() (model.Permissions, error) {
	session := r.slave
	var permissions []entity.Permission
	err := session.Order("id").Find(&permissions).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	modelPermissions := make(model.Permissions, len(permissions))
	for i, permission := range permissions {
		if p, err := permission.ConvertModel(); err != nil {
			return nil, err
		} else {
			modelPermissions[i] = *p
		}
	}
	return modelPermissions, nil
}

func (r *permission) FindByID(id model.ID) (*model.Permission, error) {
	session := r.slave
	var permission entity.Permission
	err := session.Where(&entity.Permission{Model: entity.Model{ID: *id.Value()}}).Find(&permission).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return permission.ConvertModel()
}

func (r *permission) FindByName(name model.Name) (*model.Permission, error) {
	session := r.slave
	var permission entity.Permission
	err := session.Where(&entity.Permission{Name: *name.Value()}).Find(&permission).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return permission.ConvertModel()
}

func (r *permission) Create(name model.Name, description string) (*model.Permission, error) {
	entity := entity.Permission{
		Name:        *name.Value(),
		Description: description,
	}
	entity.Generate()
	err := r.master.Create(&entity).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	return model.NewPermission(entity.ID, entity.Name, entity.Description)
}

func (r *permission) CreateBatch(names []model.Name, descriptions []string) ([]*model.Permission, error) {
	entities := make([]entity.Permission, len(names))
	for idx, name := range names {
		entity := entity.Permission{
			Name:        *name.Value(),
			Description: descriptions[idx],
		}
		entity.Generate()
		entities[idx] = entity
	}
	err := r.master.Create(&entities).Error
	if err != nil {
		return nil, model.NewDBErr(err)
	}
	mPermissions := make([]*model.Permission, len(entities))
	for idx, entity := range entities {
		mPermissions[idx], err = model.NewPermission(entity.ID, entity.Name, entity.Description)
		if err != nil {
			return nil, err
		}
	}
	return mPermissions, nil
}

func (r *permission) Update(permission *model.Permission) error {
	value := entity.Permission{
		Name:        *permission.GetName(),
		Description: permission.GetDescription(),
	}
	return model.NewDBErr(r.master.Where(&entity.Permission{Model: entity.Model{ID: *permission.GetID()}}).Updates(value).Error)
}

func (r *permission) Delete(id model.ID) error {
	db := r.master.Where(&entity.Permission{Model: entity.Model{ID: *id.Value()}}).Delete(&entity.Permission{})
	if db.RowsAffected == 0 {
		return model.ErrNoData
	}
	return model.NewDBErr(db.Error)
}
