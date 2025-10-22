package main

import (
	"fmt"
	"github.com/jukov801/Golang_MIPT/HW_4/ledger"
	"log"
	"time"
)

func main() {
	ledgerService := ledger.NewLedger()

	fmt.Println("=== Демонстрация работы интерфейса Validatable ===")

	fmt.Println("\n1. Тестирование валидной транзакции:")
	validTransaction := &ledger.Transaction{
		ID:       "1",
		Amount:   100.0,
		Category: "food",
		Date:     time.Now(),
		Type:     "expense",
	}

	err := ledger.CheckValid(validTransaction)
	if err == nil {
		err = ledgerService.AddTransaction(validTransaction)
		if err != nil {
			log.Printf("Ошибка добавления транзакции: %v", err)
		} else {
			fmt.Println("Транзакция успешно добавлена")
		}
	}

	fmt.Println("\n2. Тестирование невалидной транзакции (отрицательная сумма):")
	invalidTransaction := &ledger.Transaction{
		ID:       "2",
		Amount:   -50.0,
		Category: "food",
		Date:     time.Now(),
		Type:     "expense",
	}

	err = ledger.CheckValid(invalidTransaction)
	if err == nil {
		err = ledgerService.AddTransaction(invalidTransaction)
		if err != nil {
			log.Printf("Ошибка добавления транзакции: %v", err)
		}
	} else {
		fmt.Printf("Транзакция не добавлена: %v\n", err)
	}

	fmt.Println("\n3. Тестирование невалидной транзакции (пустая категория):")
	invalidTransaction2 := &ledger.Transaction{
		ID:       "3",
		Amount:   50.0,
		Category: "",
		Date:     time.Now(),
		Type:     "expense",
	}

	err = ledger.CheckValid(invalidTransaction2)
	if err == nil {
		err = ledgerService.AddTransaction(invalidTransaction2)
		if err != nil {
			log.Printf("Ошибка добавления транзакции: %v", err)
		}
	} else {
		fmt.Printf("Транзакция не добавлена: %v\n", err)
	}

	fmt.Println("\n4. Тестирование валидного бюджета:")
	validBudget := &ledger.Budget{
		Category: "entertainment",
		Limit:    200.0,
	}

	err = ledger.CheckValid(validBudget)
	if err == nil {
		err = ledgerService.SetBudget(validBudget.Category, validBudget.Limit)
		if err != nil {
			log.Printf("Ошибка установки бюджета: %v", err)
		} else {
			fmt.Println("Бюджет успешно установлен")
		}
	}

	fmt.Println("\n5. Тестирование невалидного бюджета (отрицательный лимит):")
	invalidBudget := &ledger.Budget{
		Category: "shopping",
		Limit:    -100.0,
	}

	err = ledger.CheckValid(invalidBudget)
	if err == nil {
		err = ledgerService.SetBudget(invalidBudget.Category, invalidBudget.Limit)
		if err != nil {
			log.Printf("Ошибка установки бюджета: %v", err)
		}
	} else {
		fmt.Printf("Бюджет не установлен: %v\n", err)
	}

	fmt.Println("\n6. Тестирование невалидного бюджета (пустая категория):")
	invalidBudget2 := &ledger.Budget{
		Category: "",
		Limit:    100.0,
	}

	err = ledger.CheckValid(invalidBudget2)
	if err == nil {
		err = ledgerService.SetBudget(invalidBudget2.Category, invalidBudget2.Limit)
		if err != nil {
			log.Printf("Ошибка установки бюджета: %v", err)
		}
	} else {
		fmt.Printf("Бюджет не установлен: %v\n", err)
	}

	fmt.Println("\n=== Демонстрация полиморфизма через интерфейс ===")

	items := []ledger.Validatable{
		&ledger.Transaction{Amount: 75.0, Category: "transport", Date: time.Now()},
		&ledger.Budget{Category: "utilities", Limit: 150.0},
		&ledger.Transaction{Amount: -25.0, Category: "food", Date: time.Now()},
		&ledger.Budget{Category: "", Limit: 300.0},
	}

	for i, item := range items {
		fmt.Printf("\nЭлемент %d:\n", i+1)
		err := ledger.CheckValid(item)
		if err != nil {
			fmt.Printf("Валидация не пройдена: %v\n", err)
		} else {
			fmt.Println("Валидация пройдена успешно")
		}
	}

	fmt.Println("\n=== Проверка состояния Ledger ===")
	fmt.Printf("Количество транзакций: %d\n", len(ledgerService.Transactions))
	fmt.Printf("Количество бюджетов: %d\n", len(ledgerService.Budgets))
}
