package service

import (
	"billingg-engine/model"
	"errors"
)

type Loan interface {
	IsDelinquent(username string) (string, error)
	GetOutStanding(user string) int64
	MakePayment(amount int64) error
}

type loan struct {
	userDataSource           []model.User
	loanDataSource           []model.Loan
	paymentHistoryDataSource []model.PaymentHistory
}

func NewLoanService(users []model.User, loans []model.Loan, histories []model.PaymentHistory) Loan {
	return loan{
		userDataSource:           users,
		loanDataSource:           loans,
		paymentHistoryDataSource: histories,
	}
}

func (l loan) GetOutStanding(username string) int64 {
	var outstanding int64

	mapLoans := make(map[string]int64)
	for _, val := range l.loanDataSource {
		// didn't include the loan that alr closed
		// and we assume that 1 cust only able to have 1 loan at the same time
		if val.Amount != 0 {
			mapLoans[val.UserID.String()] = val.Amount
		}
	}

	for _, user := range l.userDataSource {
		if val, ok := mapLoans[user.ID.String()]; ok && user.Username == username {
			outstanding = val
		}
	}

	return outstanding
}

func (l loan) IsDelinquent(username string) (string, error) {
	for _, val := range l.userDataSource {
		if username == val.Username {
			return val.Status, nil
		}
	}

	return "", errors.New("user not found")
}

func (l loan) MakePayment(amount int64) error {

	return nil
}

func Schedule() {
}
