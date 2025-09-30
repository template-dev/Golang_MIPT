package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	if username == "" {
		username = "не установлено"
	}
	fmt.Printf("Имя пользователя: %s\n", username)

	fmt.Println("\nАргументы командной строки:")
	if len(os.Args) > 1 {
		for i, arg := range os.Args[1:] {
			fmt.Printf("Аргумент %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("Аргументы не предоставлены")
	}

	fmt.Printf("\nТекущая версия Go: %s\n", runtime.Version())
}
