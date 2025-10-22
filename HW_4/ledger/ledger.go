package ledger

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID       string
	Amount   float64
	Category string
	Date     time.Time
	Type     string // "income" или "expense"
}

type Budget struct {
	Category string
	Limit    float64
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

func (l *Ledger) AddTransaction(t *Transaction) error {
	if err := t.Validate(); err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	l.Transactions = append(l.Transactions, t)
	return nil
}

func (l *Ledger) SetBudget(category string, limit float64) error {
	budget := &Budget{
		Category: category,
		Limit:    limit,
	}

	if err := budget.Validate(); err != nil {
		return fmt.Errorf("invalid budget: %w", err)
	}

	l.Budgets[category] = budget
	return nil
}
