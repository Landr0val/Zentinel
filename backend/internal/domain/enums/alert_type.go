package enums

type AlertType string

const (
	AlertTypeUnusualAmount   AlertType = "unusual_amount"
	AlertTypeUnusualLocation AlertType = "unusual_location"
	AlertTypeHighFrequency   AlertType = "high_frequency"
	AlertTypeFraudSuspicion  AlertType = "fraud_suspicion"
)
