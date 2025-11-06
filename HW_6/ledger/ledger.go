package ledger

import (
	"errors"
	"time"
)

var (
	ErrBudgetExceeded = errors.New("budget exceeded")
)

type Transaction struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
}

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Spent    float64 `json:"spent,omitempty"`
}

type Ledger struct {
	Transactions []*Transaction
	Budgets      map[string]*Budget
}

func NewLedger() *Ledger {
	return &Ledger{
		Transactions: make([]*Transaction, 0),
		Budgets:      make(map[string]*Budget),
	}
}

func (l *Ledger) AddTransaction(tx *Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	if tx.Type == "expense" {
		budget, exists := l.Budgets[tx.Category]
		if exists {
			currentSpent := l.GetCategorySpending(tx.Category)
			if currentSpent+tx.Amount > budget.Limit {
				return ErrBudgetExceeded
			}
		}
	}

	l.Transactions = append(l.Transactions, tx)
	return nil
}

func (l *Ledger) SetBudget(b *Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}

	l.Budgets[b.Category] = b
	return nil
}

func (l *Ledger) ListTransactions() []*Transaction {
	return l.Transactions
}

func (l *Ledger) ListBudgets() []*Budget {
	budgets := make([]*Budget, 0, len(l.Budgets))
	for _, budget := range l.Budgets {
		budgets = append(budgets, budget)
	}
	return budgets
}

func (l *Ledger) GetCategorySpending(category string) float64 {
	var total float64
	for _, tx := range l.Transactions {
		if tx.Category == category && tx.Type == "expense" {
			total += tx.Amount
		}
	}
	return total
}

func (l *Ledger) Reset() {
	l.Transactions = make([]*Transaction, 0)
	l.Budgets = make(map[string]*Budget)
}
