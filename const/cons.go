package _const

const (
	DefaultPaymentAmount int64 = 110000
	RoleAdmin                  = "admin"
	RoleCustomer               = "customer"
)

type contextKey string

const UserContextKey = contextKey("user")