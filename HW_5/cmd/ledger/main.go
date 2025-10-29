// cmd/ledger/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jukov801/Golang_MIPT/HW_4/ledger"
)

func main() {
	ledgerService := ledger.NewLedger()
	handler := ledger.NewHandler(ledgerService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/transactions", handler.CreateTransactionHandler)
	mux.HandleFunc("GET /api/transactions", handler.ListTransactionsHandler)

	mux.HandleFunc("POST /api/budgets", handler.CreateBudgetHandler)
	mux.HandleFunc("GET /api/budgets", handler.ListBudgetsHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	handlerWithMiddleware := ledger.LoggingMiddleware(mux)

	port := ":8080"
	fmt.Printf("Ledger server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/transactions - Create transaction")
	fmt.Println("  GET  /api/transactions - List transactions")
	fmt.Println("  POST /api/budgets      - Create budget")
	fmt.Println("  GET  /api/budgets      - List budgets")
	fmt.Println("  GET  /health           - Health check")

	log.Fatal(http.ListenAndServe(port, handlerWithMiddleware))
}
