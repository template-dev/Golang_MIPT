## Русский

### Задание. Напишите программу на Go, которая:

* выводит имя пользователя из переменной окружения;
* считывает аргументы CLI и выводит их;
* показывает текущую версию Go через `runtime.Version()`.

### Запуск программы

```
go run main.go testArg
```

### Вывод

````
Имя пользователя: jukov

Аргументы командной строки:
Аргумент 1: testArg

Текущая версия Go: go1.24.2
````

---

## English

### Task. Write a Go program that:

* prints the username from an environment variable;
* reads and prints the CLI arguments;
* shows the current Go version using `runtime.Version()`.

### Running the program

```
go run main.go testArg
```

### Output

```
Username: jukov

Command line arguments:
Argument 1: testArg

Current Go version: go1.24.2
```