package _const

type contextKey string

const (
	DefaultPaymentAmount int64 = 110000
	RoleAdmin                  = "admin"
	RoleCustomer               = "customer"
	UserContextKey             = contextKey("user")
	StatusClean                = "Clean"
	StatusDelinquent           = "Delinquent"
	TotalPaymentWeek           = 50
	TotalPaymentLoan           = 5500000

	UserNotFoundErr      = "user not found"
	InsufficientPaidErr  = "insufficient amount for pay loan"
	AlreadyPaidErr       = "loan already paid this week"
	CurrentLoanExistsErr = "you still have on going loan, can't request new loan"
)
