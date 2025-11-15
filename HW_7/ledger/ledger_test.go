package ledger

import (
	"testing"
	"time"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name        string
		transaction Transaction
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid transaction",
			transaction: Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Now().Add(-24 * time.Hour),
				Type:     "expense",
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			transaction: Transaction{
				Amount:   0,
				Category: "food",
				Date:     time.Now(),
				Type:     "expense",
			},
			wantErr: true,
			errMsg:  "amount must be positive",
		},
		{
			name: "negative amount",
			transaction: Transaction{
				Amount:   -50.0,
				Category: "food",
				Date:     time.Now(),
				Type:     "expense",
			},
			wantErr: true,
			errMsg:  "amount must be positive",
		},
		{
			name: "empty category",
			transaction: Transaction{
				Amount:   100.0,
				Category: "",
				Date:     time.Now(),
				Type:     "expense",
			},
			wantErr: true,
			errMsg:  "category cannot be empty",
		},
		{
			name: "zero date",
			transaction: Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Time{},
				Type:     "expense",
			},
			wantErr: true,
			errMsg:  "date cannot be zero",
		},
		{
			name: "future date",
			transaction: Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Now().Add(24 * time.Hour),
				Type:     "expense",
			},
			wantErr: true,
			errMsg:  "date cannot be in the future",
		},
		{
			name: "invalid type",
			transaction: Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Now(),
				Type:     "invalid",
			},
			wantErr: true,
			errMsg:  "type must be 'income' or 'expense'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.transaction.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBudget_Validate(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid budget",
			budget: Budget{
				Category: "food",
				Limit:    1000.0,
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: Budget{
				Category: "food",
				Limit:    0,
			},
			wantErr: true,
			errMsg:  "limit must be positive",
		},
		{
			name: "negative limit",
			budget: Budget{
				Category: "food",
				Limit:    -100.0,
			},
			wantErr: true,
			errMsg:  "limit must be positive",
		},
		{
			name: "empty category",
			budget: Budget{
				Category: "",
				Limit:    1000.0,
			},
			wantErr: true,
			errMsg:  "category cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.budget.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestLedger_BudgetExceeded(t *testing.T) {
	ledger := NewLedger()

	t.Cleanup(func() {
		ledger.Reset()
	})

	budget := &Budget{Category: "food", Limit: 5000.0}
	if err := ledger.SetBudget(budget); err != nil {
		t.Fatalf("Failed to set budget: %v", err)
	}

	t.Run("transaction within budget", func(t *testing.T) {
		tx := &Transaction{
			ID:       "1",
			Amount:   1000.0,
			Category: "food",
			Date:     time.Now(),
			Type:     "expense",
		}

		err := ledger.AddTransaction(tx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		transactions := ledger.ListTransactions()
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}
	})

	t.Run("transaction exceeds budget", func(t *testing.T) {
		initialCount := len(ledger.ListTransactions())

		tx := &Transaction{
			ID:       "2",
			Amount:   4500.0,
			Category: "food",
			Date:     time.Now(),
			Type:     "expense",
		}

		err := ledger.AddTransaction(tx)
		if err != ErrBudgetExceeded {
			t.Errorf("Expected ErrBudgetExceeded, got %v", err)
		}

		transactions := ledger.ListTransactions()
		if len(transactions) != initialCount {
			t.Errorf("Expected %d transactions after rejection, got %d", initialCount, len(transactions))
		}
	})

	t.Run("income transactions ignore budget", func(t *testing.T) {
		tx := &Transaction{
			ID:       "3",
			Amount:   10000.0,
			Category: "food",
			Date:     time.Now(),
			Type:     "income",
		}

		err := ledger.AddTransaction(tx)
		if err != nil {
			t.Errorf("Expected no error for income transaction, got %v", err)
		}
	})
}

func TestLedger_ListFunctions(t *testing.T) {
	ledger := NewLedger()

	t.Cleanup(func() {
		ledger.Reset()
	})

	t.Run("empty lists", func(t *testing.T) {
		transactions := ledger.ListTransactions()
		if len(transactions) != 0 {
			t.Errorf("Expected 0 transactions, got %d", len(transactions))
		}

		budgets := ledger.ListBudgets()
		if len(budgets) != 0 {
			t.Errorf("Expected 0 budgets, got %d", len(budgets))
		}
	})

	t.Run("with data", func(t *testing.T) {
		budget := &Budget{Category: "transport", Limit: 2000.0}
		ledger.SetBudget(budget)

		tx := &Transaction{
			ID:       "1",
			Amount:   500.0,
			Category: "transport",
			Date:     time.Now(),
			Type:     "expense",
		}
		ledger.AddTransaction(tx)

		transactions := ledger.ListTransactions()
		if len(transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(transactions))
		}

		budgets := ledger.ListBudgets()
		if len(budgets) != 1 {
			t.Errorf("Expected 1 budget, got %d", len(budgets))
		}
	})
}
