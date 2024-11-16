package service

import (
	"billingg-engine/model"
	"errors"
)

type Loan interface {
	IsDelinquent(username string) (string, error)
	GetOutStanding() int64
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

func (l loan) GetOutStanding() int64 {
	var outstanding int64
	for _, loan := range l.loanDataSource {
		for _, user := range l.userDataSource {
			if loan.UserID == user.ID {
				outstanding = loan.Amount
			}
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
