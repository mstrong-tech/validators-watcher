package alerts

import (
	"validators-watcher/db"
	"validators-watcher/validators"
)

type AlertsChecker interface {
	CheckAlerts(
		network string,
		actualState validators.ValidatorData,
		previousStates []db.ValidatorDBItem,
	) (string, error)
}

type AlertsSender interface {
	SendAlert(
		alert string,
	) error
}
