package service

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"billingg-engine/util"
	"errors"
	"github.com/sirupsen/logrus"
	"math"
	"time"
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
	var user model.User
	for _, val := range l.userDataSource {
		if username == val.Username {
			user = val
			break
		}
	}

	var loan model.Loan
	for _, val := range l.loanDataSource {
		if val.UserID == user.ID && val.Amount > 0 {
			loan = val
			break
		}
	}

	var histories []model.PaymentHistory
	for _, val := range l.paymentHistoryDataSource {
		if val.LoanID == loan.ID {
			histories = append(histories, val)
		}
	}

	// if user are not found
	if (user == model.User{}) {
		return "", errors.New(_const.UserNotFoundErr)
	}

	// check if their loan is new
	if l.isNewLoan(user) {
		logrus.Infof("user %v is just create his loan this week (%v)", user.Username, loan.CreatedAt.Weekday())
		return _const.StatusClean, nil
	}

	// check based on their payment histories and week passed
	{
		_, weekLoan := loan.CreatedAt.ISOWeek()
		_, currentWeek := time.Now().ISOWeek()
		weekPassed := currentWeek - weekLoan
		logrus.Info("week passed:", weekPassed)
		logrus.Info("total payment:", len(histories))

		if int(math.Abs(float64(weekPassed-len(histories)))) >= 2 {
			return _const.StatusDelinquent, nil
		}
	}

	return _const.StatusClean, nil
}

func (l loan) MakePayment(amount int64) error {

	return nil
}

func Schedule() {
}

func (l loan) isNewLoan(user model.User) bool {
	for _, val := range l.loanDataSource {
		if val.Amount > 0 && val.UserID == user.ID {
			if util.IsInCurrentWeek(val.CreatedAt) {
				return true
			}
		}
	}
	return false
}
