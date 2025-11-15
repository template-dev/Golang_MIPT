package ledger

import (
	"errors"
	"time"
)

type Validatable interface {
	Validate() error
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if t.Category == "" {
		return errors.New("category cannot be empty")
	}

	if t.Date.IsZero() {
		return errors.New("date cannot be zero")
	}

	if t.Date.After(time.Now()) {
		return errors.New("date cannot be in the future")
	}

	return nil
}

func (b *Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("limit must be positive")
	}

	if b.Category == "" {
		return errors.New("category cannot be empty")
	}

	return nil
}
