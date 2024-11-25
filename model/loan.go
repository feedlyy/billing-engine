package model

import (
	"github.com/google/uuid"
	"time"
)

type Tz struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Loan struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Amount int64 `json:"amount"`
	Tz
}

type PaymentHistory struct {
	ID     uuid.UUID
	LoanID uuid.UUID
	Amount int64
	Tz
}
