# SIS Expense Analyzer — Checkpoint 1

A synchronous, in-memory Go application for analyzing expense batches. Each batch is a **task** with its own identifier, input, lifecycle state, and stored result.

## Requirements

- Go 1.27
- No third-party packages or external services

## Run

```bash
go run .
```

## Verify

```bash
gofmt -l .
go vet ./...
go test ./...
```

`gofmt -l .` should print nothing; vet and tests should succeed.

## Project files

- `task.go`: Task, Expense, Result, states, errors, and safe cloning
- `processor.go`: Processor interface, input validation, analysis
- `store.go`: In-memory task creation, lookup, processing
- `main.go`: Demo using two independent tasks
- `task_test.go`: Automated tests, including table-driven validation cases

## Input

Each task contains one or more expenses. Each expense has:
- `Category`: nonempty string
- `Description`: nonempty string
- `AmountCents`: positive integer amount in cents (e.g., 1599 means $15.99)

Using integer cents avoids floating-point rounding problems for totals.

## Processing result

The processor computes:
- Total amount
- Average amount
- Highest individual expense
- Total amount by category
- Number of expenses

## Lifecycle

`Pending -> Processing -> Completed`

- `Create` validates input and stores a **Pending** task.
- `Process` accepts only a Pending task, sets it to **Processing** while calculations run, stores the result, then sets **Completed**.
- Reprocessing a Completed task returns `ErrInvalidState`.
- If processing fails, the task returns to Pending; the error is returned.
- `Get` and `Process` on unknown IDs return `ErrNotFound`.

Processing is synchronous: the Processing state is internal during the call.

## Design choices

- `Task`, `Expense`, and `Result` are structs.
- `Processor` is a meaningful interface so the store depends on behavior rather than a concrete calculator.
- `ExpenseProcessor` implements the interface with value receivers.
- `Store` uses pointer receivers because it modifies an in-memory `map[int]*Task` and an ID counter.
- The store clones inputs and returned tasks to prevent accidental external modification of stored data.
- Errors are returned, not panicked.
- No concurrency, HTTP, database, or external APIs (reserved for future checkpoints).
- IDs are unique within the lifetime of one Store.

## Tests

`go test ./...` covers:
1. Successful creation
2. Invalid inputs (table-driven)
3. Successful retrieval
4. Missing task retrieval
5. Successful processing
6. Result correctness
7. State after processing
8. Invalid repeated processing
9. Multiple independent tasks and safe copies

## Example

Expenses: Food/Lunch 1500 cents, Transport/Taxi 1000, Food/Dinner 2500, Shopping/Headphones 5000.

Result: total 10000 cents ($100.00), average 2500 cents ($25.00), highest Headphones 5000 cents ($50.00), categories Food 4000, Transport 1000, Shopping 5000.
