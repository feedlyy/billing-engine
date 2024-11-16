package handler

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"billingg-engine/service"
	"encoding/json"
	"net/http"
)

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
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		return
	}

	type Result struct {
		Outstanding int64 `json:"outstanding"`
	}

	outstanding := l.svc.GetOutStanding(loggedUser.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Result{Outstanding: outstanding})
}
