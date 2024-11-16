package handler

import (
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("test")
}
