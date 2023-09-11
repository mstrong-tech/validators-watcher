package db

import (
	"strings"
	"validators-watcher/validators"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDBOptions struct {
	Path string `json:"path" yaml:"path"`
}

type SqliteDB struct {
	db *gorm.DB
}

func NewSqliteDB(opts SqliteDBOptions) (ValidatorsDB, error) {
	var validatorsDB SqliteDB

	db, err := gorm.Open(sqlite.Open(opts.Path), &gorm.Config{})
	db.AutoMigrate(&NetworkModel{})
	db.AutoMigrate(&ValidatorModel{})

	validatorsDB.db = db

	return &validatorsDB, err
}

type NetworkModel struct {
	gorm.Model
	Name       string
	Validators []ValidatorModel `gorm:"foreignKey:NetworkID"`
}

type ValidatorModel struct {
	gorm.Model
	Index                 string
	Pubkey                string
	Balance               string
	EffectiveBalance      string
	WithdrawalCredentials string
	NetworkID             uint
}

func (validatorModel ValidatorModel) ToValidatorDBItem() ValidatorDBItem {
	return ValidatorDBItem{
		Index:                 validatorModel.Index,
		Pubkey:                validatorModel.Pubkey,
		Balance:               validatorModel.Balance,
		EffectiveBalance:      validatorModel.EffectiveBalance,
		WithdrawalCredentials: validatorModel.WithdrawalCredentials,
		CreatedAt:             validatorModel.CreatedAt,
	}
}

func (db *SqliteDB) Update(network string, validator validators.ValidatorData) error {
	network = strings.ToLower(network)

	var networkModel NetworkModel
	if err := db.db.FirstOrCreate(&networkModel,
		NetworkModel{
			Name: network,
		}).Error; err != nil {
		return err
	}

	var validatorModel ValidatorModel
	if err := db.db.FirstOrCreate(&validatorModel,
		ValidatorModel{
			Index:                 validator.Data.Index,
			Pubkey:                validator.Data.Validator.Pubkey,
			Balance:               validator.Data.Balance,
			EffectiveBalance:      validator.Data.Validator.EffectiveBalance,
			WithdrawalCredentials: validator.Data.Validator.WithdrawalCredentials,
			NetworkID:             networkModel.ID,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (db *SqliteDB) Get(network string, validator validators.ValidatorData) (ValidatorDBItem, error) {
	network = strings.ToLower(network)
	var networkModel NetworkModel
	if err := db.db.Model(&NetworkModel{}).Where(&NetworkModel{
		Name: network,
	}).First(&networkModel).Error; err != nil {
		return ValidatorDBItem{}, err
	}

	var validatorModel ValidatorModel
	if err := db.db.Model(&ValidatorModel{}).Order("created_at desc").First(&validatorModel, &ValidatorModel{
		Index:     validator.Data.Index,
		Pubkey:    validator.Data.Validator.Pubkey,
		NetworkID: networkModel.ID,
	}).Error; err != nil {
		return ValidatorDBItem{}, err
	}

	return validatorModel.ToValidatorDBItem(), nil
}

func (db *SqliteDB) GetLatest(
	network string,
	validator validators.ValidatorData,
	limit int,
) ([]ValidatorDBItem, error) {
	network = strings.ToLower(network)
	var networkModel NetworkModel
	if err := db.db.Model(&NetworkModel{}).Where(&NetworkModel{
		Name: network,
	}).First(&networkModel).Error; err != nil {
		return nil, err
	}

	var validatorModels []ValidatorModel
	if err := db.db.Model(&ValidatorModel{}).Order("created_at desc").Limit(limit).Find(&validatorModels, &ValidatorModel{
		Index:     validator.Data.Index,
		Pubkey:    validator.Data.Validator.Pubkey,
		NetworkID: networkModel.ID,
	}).Error; err != nil {
		return nil, err
	}

	var validatorDBItems []ValidatorDBItem
	for _, validatorModel := range validatorModels {
		validatorDBItems = append(validatorDBItems, validatorModel.ToValidatorDBItem())
	}

	return validatorDBItems, nil
}
