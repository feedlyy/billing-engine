package handler

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"billingg-engine/service"
	"billingg-engine/util"
	"net/http"
	"strconv"
)

type Err struct {
	Error string `json:"error"`
}

type LoanHandler struct {
	svc service.Loan
}

func NewLoanHandler(loan service.Loan) LoanHandler {
	return LoanHandler{
		svc: loan,
	}
}

func (l LoanHandler) GetCurrentOutStanding(w http.ResponseWriter, r *http.Request) {
	loggedUser, ok := r.Context().Value(_const.UserContextKey).(*model.UserClaims)
	if !ok {
		util.RespErr(w, Err{Error: "Failed to retrieve user information"}, http.StatusInternalServerError)
		return
	}

	type Result struct {
		Outstanding int64 `json:"outstanding"`
	}

	outstanding := l.svc.GetOutStanding(loggedUser.Username)
	util.RespOK(w, Result{Outstanding: outstanding})
}

func (l LoanHandler) CheckIsDelinquent(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	user := q.Get("user")

	if user == "" {
		util.RespErr(w, Err{Error: "Missing or empty value query param {user}"}, http.StatusBadRequest)
		return
	}

	type Result struct {
		Status string `json:"status"`
	}

	status, err := l.svc.IsDelinquent(user)
	if err != nil {
		switch err.Error() {
		case _const.UserNotFoundErr:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusNotFound)
			return
		default:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
	}
	util.RespOK(w, Result{Status: status})
}

func (l LoanHandler) Payment(w http.ResponseWriter, r *http.Request) {
	loggedUser, ok := r.Context().Value(_const.UserContextKey).(*model.UserClaims)
	if !ok {
		util.RespErr(w, Err{Error: "Failed to retrieve user information"}, http.StatusInternalServerError)
		return
	}

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		util.RespErr(w, Err{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	type Result struct{}

	err = l.svc.MakePayment(int64(amount), loggedUser.Username)
	if err != nil {
		switch err.Error() {
		case _const.UserNotFoundErr:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusNotFound)
			return
		case _const.InsufficientPaidErr:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusBadRequest)
			return
		case _const.AlreadyPaidErr:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusBadRequest)
			return
		default:
			util.RespErr(w, Err{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
	}
	util.RespOK(w, Result{})
}
