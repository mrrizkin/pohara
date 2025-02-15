package repository

import "github.com/mrrizkin/pohara/app/model"

func (a *AuthRepository) CreateCfgPolicy(cfgPolicy *model.CfgPolicy) error {
	result := a.db.Create(cfgPolicy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *AuthRepository) GetCfgPolicyByID(id uint) (*model.CfgPolicy, error) {
	var cfgPolicy model.CfgPolicy
	result := a.db.First(&cfgPolicy, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cfgPolicy, nil
}

func (a *AuthRepository) GetAllCfgPolicies() ([]model.CfgPolicy, error) {
	var cfgPolicies []model.CfgPolicy
	result := a.db.Find(&cfgPolicies)
	if result.Error != nil {
		return nil, result.Error
	}
	return cfgPolicies, nil
}

func (a *AuthRepository) UpdateCfgPolicy(cfgPolicy *model.CfgPolicy) error {
	result := a.db.Save(cfgPolicy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *AuthRepository) DeleteCfgPolicy(id uint) error {
	result := a.db.Delete(&model.CfgPolicy{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
