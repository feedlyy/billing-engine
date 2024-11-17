package util

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// initiate dummy data
var users = []model.User{
	{
		ID:           uuid.New(),
		Username:     "charlie12",
		Email:        "charlie@gmail.com",
		PasswordHash: GeneratePasswordHash("test1"),
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, -2, 0),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
			DeletedAt: time.Time{}, // Not deleted
		},
		Role: _const.RoleCustomer,
	},
	{
		ID:           uuid.New(),
		Username:     "bob23",
		Email:        "bob@gmail.com",
		PasswordHash: GeneratePasswordHash("test2"),
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, -3, 0),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
			DeletedAt: time.Time{},
		},
		Role: _const.RoleCustomer,
	},
	{
		ID:           uuid.New(),
		Username:     "steven123",
		Email:        "steven@gmail.com",
		PasswordHash: GeneratePasswordHash("test3"),
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, 0, -7),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
			DeletedAt: time.Time{},
		},
		Role: _const.RoleCustomer,
	},
	{
		ID:           uuid.New(),
		Username:     "admin123",
		Email:        "admin@gmail.com",
		PasswordHash: GeneratePasswordHash("testAdmin"),
		Tz: model.Tz{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
		Role: _const.RoleAdmin,
	},
}

var loans = []model.Loan{
	{
		ID:     uuid.New(),
		UserID: users[0].ID,
		Amount: 5390000,
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, 0, -21), // which mean is alr on the third week
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
	{
		ID:     uuid.New(),
		UserID: users[1].ID,
		Amount: 5500000,
		Tz: model.Tz{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
	{
		ID:     uuid.New(),
		UserID: users[2].ID,
		Amount: 5280000,
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, 0, -14),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
}

var paymentHistories = []model.PaymentHistory{
	{
		ID:     uuid.New(),
		LoanID: loans[0].ID,
		Amount: _const.DefaultPaymentAmount,
		Tz: model.Tz{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
	{
		ID:     uuid.New(),
		LoanID: loans[2].ID,
		Amount: _const.DefaultPaymentAmount,
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, 0, -12),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
	{
		ID:     uuid.New(),
		LoanID: loans[2].ID,
		Amount: _const.DefaultPaymentAmount,
		Tz: model.Tz{
			CreatedAt: time.Now().AddDate(0, 0, -6),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	},
}

func GeneratePasswordHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func GenerateDummyData() ([]model.User, []model.Loan, []model.PaymentHistory) {
	return users, loans, paymentHistories
}

func RespOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespErr(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func IsInCurrentWeek(input time.Time) bool {
	now := time.Now().Local()
	localInput := input.Local()

	_, currentWeek := now.ISOWeek()
	_, inputWeek := localInput.ISOWeek()

	return currentWeek == inputWeek
}
