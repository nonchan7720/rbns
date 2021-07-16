package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiKey struct {
	AccessToken string `gorm:"type:varchar(256)"`
}

func (a *ApiKey) Generate() {
	a.AccessToken = uuid.NewString()
}

func authMigration(db *gorm.DB) error {
	if err := db.AutoMigrate(&ApiKey{}); err != nil {
		return err
	}
	var count int64
	if err := db.Model(&ApiKey{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		authKey := ApiKey{}
		authKey.Generate()
		return db.Create(&authKey).Error
	}

	return nil
}

func authMigrationBack(db *gorm.DB) error {
	return db.Migrator().DropTable(&ApiKey{})
}
