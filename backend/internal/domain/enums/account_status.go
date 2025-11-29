package enums

type AccountStatus string

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusBlocked AccountStatus = "blocked"
	AccountStatusClosed  AccountStatus = "closed"
)
