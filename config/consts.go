package config

const (
	DEFAULT_FREQUENCY         = 12 // Ethereum mainnet slot time is 12 seconds
	DEFAULT_CHECK_BALANCE_FOR = 3  // FIXME: Would it be better to wait for an entire epoch?
	DEFAULT_ALERTS_SLEEP_FOR  = 60 // Time to doze off target alerts after an alert is sent
)
