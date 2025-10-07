package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}

	tx.ID = len(transactions) + 1

	if tx.Date == "" {
		tx.Date = time.Now().Format("2006-01-02")
	}

	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func main() {
	fmt.Println("Ledger service started")

	testTransactions := []Transaction{
		{
			Amount:      1500.50,
			Category:    "Salary",
			Description: "Monthly salary",
			Date:        "2024-01-15",
		},
		{
			Amount:      45.75,
			Category:    "Food",
			Description: "Groceries",
			Date:        "2024-01-16",
		},
		{
			Amount:      120.00,
			Category:    "Entertainment",
			Description: "Cinema tickets",
			Date:        "2024-01-17",
		},
	}

	for _, tx := range testTransactions {
		err := AddTransaction(tx)
		if err != nil {
			fmt.Printf("Error adding transaction: %v\n", err)
		} else {
			fmt.Printf("Transaction added successfully: %s (%.2f)\n", tx.Category, tx.Amount)
		}
	}

	fmt.Println("\nAll transactions:")
	allTransactions := ListTransactions()
	for _, tx := range allTransactions {
		fmt.Printf("ID: %d | Date: %s | Category: %s | Amount: %.2f | Description: %s\n",
			tx.ID, tx.Date, tx.Category, tx.Amount, tx.Description)
	}

	fmt.Printf("\nTotal transactions: %d\n", len(allTransactions))
}
