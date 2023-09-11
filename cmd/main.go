package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
	"validators-watcher/alerts"
	"validators-watcher/config"
	"validators-watcher/db"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	// Parse flags
	logLevel := flag.Uint("log-level", uint(logrus.InfoLevel), "log level")
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Set logging level
	logrus.SetLevel(logrus.Level(*logLevel))

	// Read config file
	data, err := os.ReadFile(*configPath)
	if err != nil {
		logrus.Fatalf("failed to read config file: %v", err)
	}

	var configs config.Config
	if strings.HasSuffix(*configPath, "json") {
		if err := json.Unmarshal(data, &configs); err != nil {
			logrus.Fatalf("failed to parse config file: %v", err)
		}
	} else if strings.HasSuffix(*configPath, "yaml") {
		if err := yaml.Unmarshal(data, &configs); err != nil {
			logrus.Fatalf("failed to parse config file: %v", err)
		}
	} else {
		logrus.Fatalf("unknown config file format")
	}

	completeConfig := configs.BuildCompleteConfig()

	// Load validators database
	var validatorsDB db.ValidatorsDB
	if validatorsDB, err = db.NewSqliteDB(completeConfig.Data.SqliteDBOpts); err != nil {
		logrus.Fatalf("failed to initialize db: %v", err)
	}

	// Initialize alerts sender
	var alertsSender alerts.AlertsSender

	// Start validators monitoring loops
	var wg sync.WaitGroup
	done := make(chan any)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		sig := <-signalChan
		logrus.Debugf("received signal %v", sig)
		close(done)
	}()

	for _, target := range completeConfig.Targets {
		wg.Add(1)
		go targetLoop(
			&wg,
			done,
			target,
			validatorsDB,
			alertsSender,
		)
	}

	wg.Wait()
}

func targetLoop(
	wg *sync.WaitGroup,
	done chan any,
	target config.ConfigTarget,
	validatorsDB db.ValidatorsDB,
	alertsSender alerts.AlertsSender,
) {
	firstRound := true
	ticker := time.Tick(time.Duration(target.Frequency) * time.Second)
	for {
		select {
		case <-done:
			wg.Done()
			return
		case <-ticker:
			for _, validator := range target.Validators {
				validatorData, err := target.BeaconApi.GetValidatorData(validator)
				if err != nil {
					logrus.Errorf("failed to get validator %s data: %v", validator.Index, err)
					continue
				}

				if !firstRound {
					_, err := validatorsDB.Get(target.BeaconApi.Network, validatorData)
					if err != nil {
						logrus.Errorf("failed to get validator %s data from db: %v", validator.Index, err)
						continue
					}

					// TODO: compare data and send alert if needed
				}

				if err := validatorsDB.Update(target.BeaconApi.Network, validatorData); err != nil {
					logrus.Errorf("failed to update db: %v", err)
					continue
				}
			}
			firstRound = false
		}
	}
}
