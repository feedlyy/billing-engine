package service

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"github.com/google/uuid"
	"testing"
)

func TestBuildUserInfo(t *testing.T) {
	userID := uuid.New()
	loanID := uuid.New()
	paymentHistoryID := uuid.New()

	tests := []struct {
		name                     string
		username                 string
		userDataSource           []model.User
		loanDataSource           []model.Loan
		paymentHistoryDataSource []model.PaymentHistory
		wantUser                 model.User
		wantLoan                 model.Loan
		wantHistories            []model.PaymentHistory
		wantErr                  bool
	}{
		{
			name:     "User Found With Loan and Payment History",
			username: "existing-user",
			userDataSource: []model.User{
				{ID: userID, Username: "existing-user"},
			},
			loanDataSource: []model.Loan{
				{ID: loanID, UserID: userID, Amount: 1000},
			},
			paymentHistoryDataSource: []model.PaymentHistory{
				{ID: paymentHistoryID, LoanID: loanID, Amount: 500},
			},
			wantUser: model.User{ID: userID, Username: "existing-user"},
			wantLoan: model.Loan{ID: loanID, UserID: userID, Amount: 1000},
			wantHistories: []model.PaymentHistory{
				{ID: paymentHistoryID, LoanID: loanID, Amount: 500},
			},
			wantErr: false,
		},
		{
			name:     "User Found Without Loan",
			username: "no-loan-user",
			userDataSource: []model.User{
				{ID: userID, Username: "no-loan-user"},
			},
			loanDataSource:           []model.Loan{},
			paymentHistoryDataSource: []model.PaymentHistory{},
			wantUser:                 model.User{ID: userID, Username: "no-loan-user"},
			wantLoan:                 model.Loan{}, // Zero value loan
			wantHistories:            nil,
			wantErr:                  false,
		},
		{
			name:                     "User Not Found",
			username:                 "non-existent-user",
			userDataSource:           []model.User{},
			loanDataSource:           []model.Loan{},
			paymentHistoryDataSource: []model.PaymentHistory{},
			wantUser:                 model.User{},
			wantLoan:                 model.Loan{},
			wantHistories:            nil,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loan{
				userDataSource:           tt.userDataSource,
				loanDataSource:           tt.loanDataSource,
				paymentHistoryDataSource: tt.paymentHistoryDataSource,
			}

			user, loan, histories, err := l.BuildUserInfo(tt.username)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if err.Error() != _const.UserNotFoundErr {
					t.Errorf("Expected error '%s', but got '%s'", _const.UserNotFoundErr, err.Error())
				}
				return
			}

			// Check user
			if user != tt.wantUser {
				t.Errorf("Expected user %v, but got %v", tt.wantUser, user)
			}

			// Check loan
			if loan != tt.wantLoan {
				t.Errorf("Expected loan %v, but got %v", tt.wantLoan, loan)
			}

			// Check payment histories
			if len(histories) != len(tt.wantHistories) {
				t.Errorf("Expected %d payment histories, but got %d", len(tt.wantHistories), len(histories))
				return
			}
			for i := 0; i < len(histories); i++ {
				if histories[i].ID != tt.wantHistories[i].ID || histories[i].LoanID != tt.wantHistories[i].LoanID || histories[i].Amount != tt.wantHistories[i].Amount {
					t.Errorf("History at index %d mismatch. Expected %+v, but got %+v", i, tt.wantHistories[i], histories[i])
				}
			}
		})
	}
}
