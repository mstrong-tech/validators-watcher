package alerts

import (
	"fmt"
	"strings"
	"validators-watcher/db"
	"validators-watcher/validators"
)

type BalanceAlertsChecker struct {
}

func (gac BalanceAlertsChecker) CheckAlerts(
	network string,
	actualState validators.ValidatorData,
	previousStates []db.ValidatorDBItem,
) (string, error) {
	// Check if withdrawals are happening
	_ = strings.HasPrefix(
		actualState.Data.Validator.WithdrawalCredentials,
		"0x01",
	)

	// Check if balance is decreasing for all previous states
	balanceDecreasing := true
	for i, previousState := range previousStates {
		if i == 0 {
			balanceDecreasing = actualState.Data.Balance < previousState.Balance
		} else {
			balanceDecreasing = balanceDecreasing && (previousState.Balance < previousStates[i-1].Balance)
		}
	}
	if balanceDecreasing {
		baseAlert := fmt.Sprintf(
			"ALERT: %s Validator %s balance have been decreasing for %d consecutive slots.",
			strings.ToTitle(network),
			actualState.Data.Index,
			len(previousStates),
		)
		return baseAlert, nil
	}

	return "", nil
}
