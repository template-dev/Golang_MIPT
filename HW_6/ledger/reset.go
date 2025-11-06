package ledger

func (l *Ledger) Reset() {
	l.Transactions = make([]*Transaction, 0)
	l.Budgets = make(map[string]*Budget)
}
