package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jukov801/ledger/ledger"
)

type Handler struct {
	ledger *ledger.Ledger
}

func NewHandler(ledger *ledger.Ledger) *Handler {
	return &Handler{ledger: ledger}
}

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
	Date        string  `json:"date"` // ISO format "2006-01-02"
	Type        string  `json:"type"` // "income" or "expense"
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Spent    float64 `json:"spent"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	tx := &ledger.Transaction{
		ID:          uuid.New().String(),
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        date,
		Type:        req.Type,
	}

	if err := h.ledger.AddTransaction(tx); err != nil {
		switch err {
		case ledger.ErrBudgetExceeded:
			writeError(w, http.StatusConflict, "budget exceeded")
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	response := TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
		Type:        tx.Type,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	transactions := h.ledger.ListTransactions()
	response := make([]TransactionResponse, len(transactions))

	for i, tx := range transactions {
		response[i] = TransactionResponse{
			ID:          tx.ID,
			Amount:      tx.Amount,
			Category:    tx.Category,
			Description: tx.Description,
			Date:        tx.Date,
			Type:        tx.Type,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	budget := &ledger.Budget{
		Category: req.Category,
		Limit:    req.Limit,
	}

	if err := h.ledger.SetBudget(budget); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	spent := h.ledger.GetCategorySpending(req.Category)
	response := BudgetResponse{
		Category: budget.Category,
		Limit:    budget.Limit,
		Spent:    spent,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) ListBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	budgets := h.ledger.ListBudgets()
	response := make([]BudgetResponse, len(budgets))

	for i, budget := range budgets {
		spent := h.ledger.GetCategorySpending(budget.Category)
		response[i] = BudgetResponse{
			Category: budget.Category,
			Limit:    budget.Limit,
			Spent:    spent,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
