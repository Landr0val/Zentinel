package enums

type AlertStatus string

const (
	AlertStatusPending   AlertStatus = "pending"
	AlertStatusReviewed  AlertStatus = "reviewed"
	AlertStatusDismissed AlertStatus = "dismissed"
	AlertStatusEscalated AlertStatus = "escalated"
)
