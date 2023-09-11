package alerts

import (
	"github.com/sirupsen/logrus"
)

type SimpleAlertSender struct {
}

func (s SimpleAlertSender) SendAlert(alert string) error {
	logrus.Warn(alert)
	return nil
}
