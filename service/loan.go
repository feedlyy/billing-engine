package service

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"billingg-engine/util"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

type Loan interface {
	IsDelinquent(username string) (string, error)
	GetOutStanding(user string) int64
	MakePayment(amount int64, username string) error
	Schedule(username string) []string
	Create(username string) error
}

type loan struct {
	userDataSource           []model.User
	loanDataSource           []model.Loan
	paymentHistoryDataSource []model.PaymentHistory
}

func NewLoanService(users []model.User, loans []model.Loan, histories []model.PaymentHistory) Loan {
	return &loan{
		userDataSource:           users,
		loanDataSource:           loans,
		paymentHistoryDataSource: histories,
	}
}

func (l *loan) GetOutStanding(username string) int64 {
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

func (l *loan) IsDelinquent(username string) (string, error) {
	user, loan, histories, err := l.BuildUserInfo(username)
	if err != nil {
		return "", err
	}

	// check if their loan is new
	if l.isNewLoan(user) {
		logrus.Infof("user %v is just create his loan this week (%v)", user.Username, loan.CreatedAt.Weekday())
		return _const.StatusClean, nil
	}

	// check based on their payment histories and week passed
	{
		// this weekLoan added 1 week due to first week of created, the system didn't charge the loan to be paid
		_, weekLoan := loan.CreatedAt.AddDate(0, 0, 7).ISOWeek()
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

func (l *loan) MakePayment(amount int64, username string) error {
	_, loan, _, err := l.BuildUserInfo(username)
	if err != nil {
		return err
	}

	if amount < _const.DefaultPaymentAmount {
		return errors.New(_const.InsufficientPaidErr)
	}

	// check if already paid
	for _, val := range l.paymentHistoryDataSource {
		if val.LoanID == loan.ID && util.IsInCurrentWeek(val.CreatedAt) {
			return errors.New(_const.AlreadyPaidErr)
		}
	}

	{
		// create history
		l.paymentHistoryDataSource = append(l.paymentHistoryDataSource, model.PaymentHistory{
			ID:     uuid.New(),
			LoanID: loan.ID,
			Amount: _const.DefaultPaymentAmount,
			Tz: model.Tz{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: time.Time{},
			},
		})

		// update loan
		for index, val := range l.loanDataSource {
			if val.ID == loan.ID {
				l.loanDataSource[index].Amount -= _const.DefaultPaymentAmount
			}
		}
	}

	return nil
}

func (l *loan) Schedule(username string) []string {
	user, loan, histories, _ := l.BuildUserInfo(username)

	if len(histories) == _const.TotalPaymentWeek {
		return []string{"You didn't have any schedule payment loan"}
	}

	var schedules []string
	var latestHistory = time.Now() // set default value
	{
		_, weekLoan := loan.CreatedAt.AddDate(0, 0, 7).ISOWeek() // add 2 weeks to exclude current week
		_, currentWeek := time.Now().ISOWeek()
		weekPassed := currentWeek - weekLoan
		totalMissed := int(math.Abs(float64(weekPassed - len(histories))))

		if len(histories) != 0 {
			latestHistory = histories[len(histories)-1].CreatedAt.AddDate(0, 0, 7) // add 1 week to get next week payment
		}
		for i := len(histories) + 1; i <= _const.TotalPaymentWeek; i++ {
			start, end := util.GetCurrentWeek(latestHistory)
			schedule := fmt.Sprintf("W%d : %v (%v - %v)", i, _const.DefaultPaymentAmount, start.Format("02 January 2006"), end.Format("02 January 2006"))

			// update each week
			latestHistory = latestHistory.AddDate(0, 0, 7)

			if totalMissed-1 >= (_const.TotalPaymentWeek-i) && !l.isNewLoan(user) {
				repayment := fmt.Sprintf("%v [Repayment]", schedule)
				schedule = repayment
			}
			schedules = append(schedules, schedule)
		}
	}

	return schedules
}

func (l *loan) Create(username string) error {
	user, loan, _, _ := l.BuildUserInfo(username)

	if loan.Amount > 0 {
		return errors.New(_const.CurrentLoanExistsErr)
	}

	l.loanDataSource = append(l.loanDataSource, model.Loan{
		ID:     uuid.New(),
		UserID: user.ID,
		Amount: _const.TotalPaymentLoan,
		Tz: model.Tz{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	})

	return nil
}

func (l *loan) isNewLoan(user model.User) bool {
	for _, val := range l.loanDataSource {
		if val.Amount > 0 && val.UserID == user.ID {
			if util.IsInCurrentWeek(val.CreatedAt) {
				return true
			}
		}
	}
	return false
}

func (l *loan) BuildUserInfo(username string) (model.User, model.Loan, []model.PaymentHistory, error) {
	var user model.User
	for _, val := range l.userDataSource {
		if username == val.Username {
			user = val
			break
		}
	}

	// if user are not found
	if (user == model.User{}) {
		return model.User{}, model.Loan{}, nil, errors.New(_const.UserNotFoundErr)
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

	return user, loan, histories, nil
}
