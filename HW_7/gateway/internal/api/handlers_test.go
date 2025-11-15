package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jukov801/Golang_MIPT/HW_6/ledger"
)

func TestBudgetHandlers(t *testing.T) {
	ledgerService := ledger.NewLedger()
	handler := NewHandler(ledgerService)

	t.Cleanup(func() {
		ledgerService.Reset()
	})

	t.Run("create valid budget", func(t *testing.T) {
		reqBody := `{"category":"food","limit":5000}`
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateBudgetHandler(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}

		contentType := rr.Header().Get("Content-Type")
		expectedContentType := "application/json; charset=utf-8"
		if contentType != expectedContentType {
			t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
		}

		var response ledger.BudgetResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.Category != "food" {
			t.Errorf("Expected category 'food', got '%s'", response.Category)
		}
		if response.Limit != 5000 {
			t.Errorf("Expected limit 5000, got %f", response.Limit)
		}
	})

	t.Run("create budget with invalid JSON", func(t *testing.T) {
		reqBody := `{"category":"food","limit":"invalid"}`
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateBudgetHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}

		var errorResp map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResp); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if _, exists := errorResp["error"]; !exists {
			t.Error("Expected error field in response")
		}
	})

	t.Run("create budget with negative limit", func(t *testing.T) {
		reqBody := `{"category":"food","limit":-100}`
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateBudgetHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	})

	t.Run("list budgets", func(t *testing.T) {
		createReqBody := `{"category":"entertainment","limit":3000}`
		createReq := httptest.NewRequest("POST", "/api/budgets", bytes.NewBufferString(createReqBody))
		createReq.Header.Set("Content-Type", "application/json")
		createRR := httptest.NewRecorder()
		handler.CreateBudgetHandler(createRR, createReq)

		listReq := httptest.NewRequest("GET", "/api/budgets", nil)
		listRR := httptest.NewRecorder()
		handler.ListBudgetsHandler(listRR, listReq)

		if status := listRR.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}

		var budgets []ledger.BudgetResponse
		if err := json.Unmarshal(listRR.Body.Bytes(), &budgets); err != nil {
			t.Fatalf("Failed to parse budgets list: %v", err)
		}

		if len(budgets) != 2 {
			t.Errorf("Expected 2 budgets, got %d", len(budgets))
		}
	})
}

func TestTransactionHandlers(t *testing.T) {
	ledgerService := ledger.NewLedger()
	handler := NewHandler(ledgerService)

	t.Cleanup(func() {
		ledgerService.Reset()
	})

	t.Run("setup budget", func(t *testing.T) {
		reqBody := `{"category":"food","limit":5000}`
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateBudgetHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Failed to setup budget: %d", rr.Code)
		}
	})

	t.Run("create valid transaction", func(t *testing.T) {
		reqBody := `{
			"amount": 1000,
			"category": "food",
			"description": "groceries",
			"date": "2024-01-15",
			"type": "expense"
		}`

		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateTransactionHandler(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}

		var response ledger.TransactionResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.Amount != 1000 {
			t.Errorf("Expected amount 1000, got %f", response.Amount)
		}
		if response.Category != "food" {
			t.Errorf("Expected category 'food', got '%s'", response.Category)
		}
		if response.Description != "groceries" {
			t.Errorf("Expected description 'groceries', got '%s'", response.Description)
		}
	})

	t.Run("transaction exceeds budget", func(t *testing.T) {
		reqBody := `{
			"amount": 4500,
			"category": "food",
			"description": "expensive dinner",
			"date": "2024-01-16",
			"type": "expense"
		}`

		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateTransactionHandler(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, status)
		}

		var errorResp map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResp); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if errorResp["error"] != "budget exceeded" {
			t.Errorf("Expected error 'budget exceeded', got '%s'", errorResp["error"])
		}
	})

	t.Run("list transactions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/transactions", nil)
		rr := httptest.NewRecorder()
		handler.ListTransactionsHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}

		var transactions []ledger.TransactionResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &transactions); err != nil {
			t.Fatalf("Failed to parse transactions list: %v", err)
		}

		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
	})

	t.Run("invalid transaction JSON", func(t *testing.T) {
		reqBody := `{
			"amount": "invalid",
			"category": "food",
			"date": "2024-01-15",
			"type": "expense"
		}`

		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateTransactionHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	})

	t.Run("transaction with invalid date format", func(t *testing.T) {
		reqBody := `{
			"amount": 1000,
			"category": "food",
			"date": "15-01-2024",
			"type": "expense"
		}`

		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.CreateTransactionHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	})
}

func TestMethodNotAllowed(t *testing.T) {
	ledgerService := ledger.NewLedger()
	handler := NewHandler(ledgerService)

	tests := []struct {
		method string
		path   string
	}{
		{"PUT", "/api/transactions"},
		{"DELETE", "/api/transactions"},
		{"PUT", "/api/budgets"},
		{"DELETE", "/api/budgets"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.method, tt.path), func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			switch tt.path {
			case "/api/transactions":
				handler.CreateTransactionHandler(rr, req)
			case "/api/budgets":
				handler.CreateBudgetHandler(rr, req)
			}

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s %s, got %d",
					http.StatusMethodNotAllowed, tt.method, tt.path, status)
			}
		})
	}
}
