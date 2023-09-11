package alerts

type AlertsSender interface {
	SendAlert(alert string) error
}
