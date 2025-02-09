package repository

import (
	"errors"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/database"
)

type SettingRepository struct {
	db *database.Database
}

type SettingRepositoryDependencies struct {
	fx.In

	Database *database.Database
}

func NewSettingRepository(deps SettingRepositoryDependencies) *SettingRepository {
	return &SettingRepository{
		db: deps.Database,
	}
}

func (r *SettingRepository) Create(setting *model.CfgSetting) error {
	isSet, err := r.IsSet()
	if err != nil {
		return err
	}

	if isSet {
		return errors.New("you already create a setting, just update it bro..")
	}

	return r.db.FirstOrCreate(setting).Error
}

func (r *SettingRepository) IsSet() (bool, error) {
	if _, err := r.Get(); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (r *SettingRepository) Get() (*model.CfgSetting, error) {
	var setting model.CfgSetting
	err := r.db.First(&setting).Error
	return &setting, err
}

func (r *SettingRepository) Update(setting *model.CfgSetting) error {
	return r.db.Save(setting).Error
}
