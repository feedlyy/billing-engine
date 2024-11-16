package util

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"time"
)

// initiate dummy data
var users = []model.User{
	{
		ID:           uuid.New(),
		Username:     "charlie12",
		Email:        "charlie@gmail.com",
		Status:       "Delinquent",
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
		Status:       "Clean",
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
		Status:       "Clean",
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
		Status:       "Clean",
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
			UpdatedAt: time.Time{},
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
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=19$m=65536,t=4,p=1$%s$%s", b64Salt, b64Hash)
}

func GenerateDummyData() ([]model.User, []model.Loan, []model.PaymentHistory) {
	return users, loans, paymentHistories
}
