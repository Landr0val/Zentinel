package enums

type OperationType string

const (
	OperationTypePurchase   OperationType = "purchase"
	OperationTypeWithdrawal OperationType = "withdrawal"
	OperationTypeTransfer   OperationType = "transfer"
	OperationTypeDeposit    OperationType = "deposit"
)
