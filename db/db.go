package db

import (
	"time"
	"validators-watcher/validators"
)

type ValidatorDBItem struct {
	Index                 string    `json:"index" yaml:"index"`
	Pubkey                string    `json:"pubkey" yaml:"pubkey"`
	Balance               string    `json:"balance" yaml:"balance"`
	EffectiveBalance      string    `json:"effective_balance" yaml:"effective_balance"`
	WithdrawalCredentials string    `json:"withdrawal_credential" yaml:"withdrawal_credential"`
	CreatedAt             time.Time `json:"created_at" yaml:"created_at"`
}

type NetworkDBItem struct {
	Name       string                     `json:"name" yaml:"name"`
	Validators map[string]ValidatorDBItem `json:"validators" yaml:"validators"`
}

type DBData struct {
	Networks map[string]NetworkDBItem `json:"networks" yaml:"networks"`
}

type ValidatorsDB interface {
	Update(network string, validator validators.ValidatorData) error
	Get(network string, validator validators.ValidatorData) (ValidatorDBItem, error)
	GetLatest(network string, validator validators.ValidatorData, limit int) ([]ValidatorDBItem, error)
}
