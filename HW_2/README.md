# Тестовое задание: Gateway и Ledger сервисы

## Задание 1 (Gateway)

**Цель:** Создать два основных сервиса и обеспечить их базовую функциональность.


*   Создайте два проекта на Go и напишите две программы: **Gateway** (HTTP-шлюз) и **Ledger** (бизнес-логика). Инициализируйте `go.mod` для каждого сервиса. Разместите их в одном репозитории (например, в отдельных папках `gateway` и `ledger`).


*   В сервисе **Gateway** напишите функцию `main`, которая запускает простой HTTP-сервер на порту (например, 8080). Реализуйте обработчик для пути `/ping`, который при GET-запросе возвращает строку `"pong"` со статусом `200 OK`. Это поможет проверить, что сервер запущен и отвечает.


*   В сервисе **Ledger** напишите функцию `main`, которая выводит в консоль сообщение `"Ledger service started"`. Пока можно не поднимать сервер, а просто убедиться, что модуль запускается (в дальнейшем здесь будет gRPC-сервер).


*   Запустите оба сервиса (**Gateway** и **Ledger**) отдельно. Убедитесь, что сервис **Gateway** доступен: отправьте HTTP-запрос на `GET /ping` (например, через `curl` или браузер) и получите ответ `"pong"`.

## Задание 2 (Ledger)

**Цель:** Реализовать модель данных и логику управления транзакциями в сервисе Ledger.

1.  **Структура Transaction:** Создайте структуру `Transaction` для представления финансовой транзакции. Включите следующие поля:
    *   `ID` — уникальный идентификатор (тип `int`), можно заполнять автоинкрементом (например, `len(slice)+1`).
    *   `Amount` — сумма (тип `float64` или `int`).
    *   `Category` — категория траты (тип `string`).
    *   `Description` — описание (тип `string`).
    *   `Date` — дата транзакции (можно использовать `string` для простоты или `time.Time`).


2.  **Хранилище в памяти:** Создайте в пакете **Ledger** хранилище транзакций в памяти (например, срез `[]Transaction`). Инициализируйте пустой срез при запуске сервиса.


3.  **Функция AddTransaction:** Реализуйте функцию `AddTransaction(tx Transaction) error`, которая добавляет переданную транзакцию в хранилище (в срез). Функция должна возвращать ошибку, если, например, сумма транзакции равна 0. При успешном добавлении возвращайте `nil`.


4.  **Функция ListTransactions:** Реализуйте функцию `ListTransactions() []Transaction`, которая возвращает все сохранённые транзакции (или их копию) из хранилища. Эта функция пригодится для вывода списка транзакций.


5.  **Тестирование логики:** В функции `main` сервиса **Ledger** вызовите `AddTransaction` несколько раз для теста: добавьте 2–3 произвольные транзакции (с разными категориями и суммами). Затем вызовите `ListTransactions()` и выведите список транзакций в консоль (через `fmt.Println` или форматированный вывод). Убедитесь, что все добавленные транзакции отображаются корректно.


6.  **Фиксация изменений:** Закоммитьте изменения и обновите репозиторий. Проследите, чтобы в репозитории были инструкции по запуску сервисов для проверки (можно обновить `README.md`).

---

# Coding Task: Gateway and Ledger Services

## Task 1 (Gateway)

**Objective:** Create two core services and ensure their basic functionality.

*   Create two Go projects and write two programs: **Gateway** (HTTP gateway) and **Ledger** (business logic). Initialize `go.mod` for each service. Place them in a single repository (e.g., in separate folders `gateway` and `ledger`).


*   In the **Gateway** service, write a `main` function that starts a simple HTTP server on a port (e.g., 8080). Implement a handler for the `/ping` path, which returns the string `"pong"` with a `200 OK` status upon a GET request. This will help verify that the server is running and responsive.


*   In the **Ledger** service, write a `main` function that prints a message `"Ledger service started"` to the console. For now, it's not necessary to start a server; just ensure the module runs (a gRPC server will be implemented here later).


*   Run both services (**Gateway** and **Ledger**) separately. Ensure the **Gateway** service is accessible: send an HTTP GET request to `/ping` (e.g., via `curl` or a browser) and receive a `"pong"` response.

## Task 2 (Ledger)

**Objective:** Implement the data model and transaction management logic in the Ledger service.

1.  **Transaction Struct:** Create a `Transaction` struct to represent a financial transaction. Include the following fields:
    *   `ID` — a unique identifier (type `int`), which can be auto-incremented (e.g., `len(slice)+1`).
    *   `Amount` — the transaction amount (type `float64` or `int`).
    *   `Category` — the spending category (type `string`).
    *   `Description` — a description (type `string`).
    *   `Date` — the transaction date (you can use `string` for simplicity or `time.Time`).


2.  **In-Memory Storage:** Create an in-memory storage for transactions (e.g., a slice `[]Transaction`) in the **Ledger** package. Initialize an empty slice when the service starts.


3.  **AddTransaction Function:** Implement the `AddTransaction(tx Transaction) error` function, which adds the provided transaction to the storage (the slice). The function should return an error if, for example, the transaction amount is 0. If the addition is successful, return `nil`.


4.  **ListTransactions Function:** Implement the `ListTransactions() []Transaction` function, which returns all stored transactions (or a copy of them) from the storage. This function will be useful for displaying the list of transactions.


5.  **Testing Logic:** In the `main` function of the **Ledger** service, call `AddTransaction` several times for testing: add 2–3 arbitrary transactions (with different categories and amounts). Then call `ListTransactions()` and print the list of transactions to the console (using `fmt.Println` or formatted output). Ensure that all added transactions are displayed correctly.


6.  **Commit Changes:** Commit the changes and update the repository. Make sure the repository includes instructions for running the services for verification (you can update the `README.md`).