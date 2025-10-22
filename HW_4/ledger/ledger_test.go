package ledger

import (
	"testing"
	"time"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name        string
		transaction *Transaction
		wantErr     bool
	}{
		{
			name: "valid transaction",
			transaction: &Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			transaction: &Transaction{
				Amount:   0,
				Category: "food",
				Date:     time.Now(),
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			transaction: &Transaction{
				Amount:   -50.0,
				Category: "food",
				Date:     time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty category",
			transaction: &Transaction{
				Amount:   100.0,
				Category: "",
				Date:     time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero date",
			transaction: &Transaction{
				Amount:   100.0,
				Category: "food",
				Date:     time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBudget_Validate(t *testing.T) {
	tests := []struct {
		name    string
		budget  *Budget
		wantErr bool
	}{
		{
			name: "valid budget",
			budget: &Budget{
				Category: "food",
				Limit:    500.0,
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: &Budget{
				Category: "food",
				Limit:    0,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			budget: &Budget{
				Category: "food",
				Limit:    -100.0,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			budget: &Budget{
				Category: "",
				Limit:    500.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.budget.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Budget.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckValid(t *testing.T) {
	validTransaction := &Transaction{
		Amount:   100.0,
		Category: "food",
		Date:     time.Now(),
	}

	invalidTransaction := &Transaction{
		Amount:   -50.0,
		Category: "food",
		Date:     time.Now(),
	}

	validBudget := &Budget{
		Category: "food",
		Limit:    500.0,
	}

	invalidBudget := &Budget{
		Category: "",
		Limit:    500.0,
	}

	if err := CheckValid(validTransaction); err != nil {
		t.Errorf("CheckValid(validTransaction) should not return error, got: %v", err)
	}

	if err := CheckValid(validBudget); err != nil {
		t.Errorf("CheckValid(validBudget) should not return error, got: %v", err)
	}

	if err := CheckValid(invalidTransaction); err == nil {
		t.Error("CheckValid(invalidTransaction) should return error")
	}

	if err := CheckValid(invalidBudget); err == nil {
		t.Error("CheckValid(invalidBudget) should return error")
	}
}
