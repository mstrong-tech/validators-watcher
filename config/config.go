package config

import (
	"fmt"
	"validators-watcher/db"
	"validators-watcher/validators"

	"github.com/sirupsen/logrus"
)

type ConfigData struct {
	SqliteDBOpts db.SqliteDBOptions `json:"sqlite_db_opts" yaml:"sqlite_db_opts"`
}

type ConfigRange struct {
	StartIndex int `json:"start" yaml:"start"`
	EndIndex   int `json:"end" yaml:"end"`
}

type ConfigTarget struct {
	BeaconApi  validators.BeaconApi   `json:"beacon_api" yaml:"beacon_api"`
	Validators []validators.Validator `json:"validators" yaml:"validators"`
	Frequency  int                    `json:"frequency,omitempty" yaml:"frequency,omitempty"`
	Ranges     []ConfigRange          `json:"ranges,omitempty" yaml:"ranges,omitempty"`
}

type ConfigAlerts struct {
	SleepAlertsFor int `json:"sleep_alerts_for,omitempty" yaml:"sleep_alerts_for,omitempty"`
}

type Config struct {
	Data    ConfigData     `json:"data" yaml:"data"`
	Alerts  ConfigAlerts   `json:"alerts" yaml:"alerts"`
	Targets []ConfigTarget `json:"targets" yaml:"targets"`
}

func (config Config) BuildCompleteConfig() Config {
	var completeConfig Config
	completeConfig.Data = config.Data
	completeConfig.Targets = make([]ConfigTarget, 0, len(config.Targets))
	for i, target := range config.Targets {
		logrus.Debugf("expanding validators ranges for target %d", i)
		if target.Frequency <= 0 {
			config.Targets[i].Frequency = 12 // Ethereum mainnet slot time
		}

		if len(target.Ranges) > 0 {
			for _, rangeItem := range target.Ranges {
				if rangeItem.StartIndex < 0 {
					logrus.Debugf("invalid start index. skipping...")
					continue
				}
				if rangeItem.EndIndex < 0 {
					logrus.Debugf("invalid end index. skipping...")
					continue
				}
				if rangeItem.StartIndex > rangeItem.EndIndex {
					logrus.Debugf("start index is greater than end index. skipping...")
					continue
				}
				for j := rangeItem.StartIndex; j <= rangeItem.EndIndex; j++ {
					config.Targets[i].Validators = append(config.Targets[i].Validators,
						validators.Validator{
							Index: fmt.Sprintf("%d", j),
						},
					)
				}
			}
		}
		completeConfig.Targets = append(completeConfig.Targets, config.Targets[i])
	}

	return completeConfig
}
