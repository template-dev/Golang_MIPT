package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Transaction struct {
	ID       int
	Amount   float64
	Category string
	Date     time.Time
	Type     string // "income" или "expense"
}

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

type Ledger struct {
	transactions []Transaction
	nextID       int
	budgets      map[string]Budget
}

func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
		nextID:       1,
		budgets:      make(map[string]Budget),
	}
}

func (l *Ledger) AddTransaction(tx Transaction) error {
	if tx.Type == "expense" {
		if budget, exists := l.budgets[tx.Category]; exists {
			currentSpent := l.getCurrentSpending(tx.Category)
			if currentSpent+tx.Amount > budget.Limit {
				return fmt.Errorf("budget exceeded for category '%s': limit %.2f, would be %.2f",
					tx.Category, budget.Limit, currentSpent+tx.Amount)
			}
		}
	}

	tx.ID = l.nextID
	l.nextID++
	l.transactions = append(l.transactions, tx)
	return nil
}

func (l *Ledger) getCurrentSpending(category string) float64 {
	var total float64
	for _, tx := range l.transactions {
		if tx.Category == category && tx.Type == "expense" {
			total += tx.Amount
		}
	}
	return total
}

func (l *Ledger) SetBudget(b Budget) {
	l.budgets[b.Category] = b
}

func (l *Ledger) GetBudget(category string) (Budget, bool) {
	budget, exists := l.budgets[category]
	return budget, exists
}

func (l *Ledger) GetTransactions() []Transaction {
	return l.transactions
}

func (l *Ledger) GetBudgets() map[string]Budget {
	return l.budgets
}

func (l *Ledger) LoadBudgets(r io.Reader) error {
	var budgets []Budget
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&budgets); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	for _, budget := range budgets {
		l.SetBudget(budget)
	}

	return nil
}

func (l *Ledger) PrintStatus() {
	fmt.Println("\n=== Текущее состояние ===")

	fmt.Println("\nБюджеты:")
	for category, budget := range l.budgets {
		spent := l.getCurrentSpending(category)
		remaining := budget.Limit - spent
		fmt.Printf("  %s: лимит %.2f, потрачено %.2f, осталось %.2f\n",
			category, budget.Limit, spent, remaining)
	}

	fmt.Println("\nТранзакции:")
	for _, tx := range l.transactions {
		fmt.Printf("  ID: %d, Сумма: %.2f, Категория: %s, Тип: %s, Дата: %s\n",
			tx.ID, tx.Amount, tx.Category, tx.Type, tx.Date.Format("2006-01-02"))
	}
}

func main() {
	ledger := NewLedger()

	fmt.Println("Установка начальных бюджетов...")
	ledger.SetBudget(Budget{Category: "еда", Limit: 5000, Period: "month"})
	ledger.SetBudget(Budget{Category: "транспорт", Limit: 2000, Period: "month"})
	ledger.SetBudget(Budget{Category: "развлечения", Limit: 3000, Period: "month"})

	if err := loadBudgetsFromFile(ledger, "budgets.json"); err != nil {
		fmt.Printf("Предупреждение: %v\n", err)
	} else {
		fmt.Println("Бюджеты успешно загружены из файла")
	}

	fmt.Println("\n=== Тестовые сценарии ===")

	fmt.Println("\n1. Добавление транзакции в пределах бюджета:")
	tx1 := Transaction{
		Amount:   1000,
		Category: "еда",
		Date:     time.Now(),
		Type:     "expense",
	}
	if err := ledger.AddTransaction(tx1); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	fmt.Println("\n2. Добавление второй транзакции:")
	tx2 := Transaction{
		Amount:   2000,
		Category: "еда",
		Date:     time.Now(),
		Type:     "expense",
	}
	if err := ledger.AddTransaction(tx2); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	fmt.Println("\n3. Попытка превысить бюджет:")
	tx3 := Transaction{
		Amount:   2500,
		Category: "еда",
		Date:     time.Now(),
		Type:     "expense",
	}
	if err := ledger.AddTransaction(tx3); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	fmt.Println("\n4. Транзакция в категории без бюджета:")
	tx4 := Transaction{
		Amount:   1000,
		Category: "образование",
		Date:     time.Now(),
		Type:     "expense",
	}
	if err := ledger.AddTransaction(tx4); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	fmt.Println("\n5. Добавление доходной транзакции:")
	tx5 := Transaction{
		Amount:   15000,
		Category: "зарплата",
		Date:     time.Now(),
		Type:     "income",
	}
	if err := ledger.AddTransaction(tx5); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	ledger.PrintStatus()

	fmt.Println("\n=== Обновление бюджета ===")
	ledger.SetBudget(Budget{Category: "еда", Limit: 6000, Period: "month"})
	fmt.Println("Лимит на еду увеличен до 6000")

	fmt.Println("\nПовторная попытка добавить транзакцию после увеличения лимита:")
	if err := ledger.AddTransaction(tx3); err != nil {
		fmt.Printf("   Ошибка: %v\n", err)
	} else {
		fmt.Println("   Транзакция успешно добавлена")
	}

	ledger.PrintStatus()
}

func loadBudgetsFromFile(ledger *Ledger, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %v", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return ledger.LoadBudgets(reader)
}
