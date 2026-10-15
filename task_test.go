package main

import (
	"errors"
	"testing"
)

func sampleExpenses() []Expense {
	return []Expense{
		{Category: "Food", Description: "Lunch", AmountCents: 1500},
		{Category: "Transport", Description: "Taxi", AmountCents: 1000},
		{Category: "Food", Description: "Dinner", AmountCents: 2500},
		{Category: "Shopping", Description: "Headphones", AmountCents: 5000},
	}
}

func TestCreateAndGet(t *testing.T) {
	store := NewStore(ExpenseProcessor{})
	created, err := store.Create(sampleExpenses())
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Status != StatusPending || created.Result != nil {
		t.Fatalf("unexpected created task: %+v", created)
	}
	got, err := store.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID || len(got.Expenses) != 4 {
		t.Fatalf("unexpected retrieved task: %+v", got)
	}
}

func TestInvalidInputs(t *testing.T) {
	cases := []struct {
		name     string
		expenses []Expense
	}{
		{"empty", nil},
		{"missing category", []Expense{{Description: "Lunch", AmountCents: 100}}},
		{"missing description", []Expense{{Category: "Food", AmountCents: 100}}},
		{"zero amount", []Expense{{Category: "Food", Description: "Lunch", AmountCents: 0}}},
		{"negative amount", []Expense{{Category: "Food", Description: "Lunch", AmountCents: -10}}},
		{"overflow", []Expense{
			{Category: "A", Description: "First", AmountCents: 9223372036854775807},
			{Category: "B", Description: "Second", AmountCents: 1},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := NewStore(ExpenseProcessor{})
			_, err := store.Create(tc.expenses)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("expected invalid input error, got %v", err)
			}
			if len(store.tasks) != 0 {
				t.Fatal("invalid task was stored")
			}
		})
	}
}

func TestMissingTask(t *testing.T) {
	store := NewStore(ExpenseProcessor{})
	if _, err := store.Get(99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get: expected not found, got %v", err)
	}
	if _, err := store.Process(99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Process: expected not found, got %v", err)
	}
}

func TestProcessResultAndState(t *testing.T) {
	store := NewStore(ExpenseProcessor{})
	task, err := store.Create(sampleExpenses())
	if err != nil {
		t.Fatal(err)
	}
	done, err := store.Process(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if done.Status != StatusCompleted || done.Result == nil {
		t.Fatalf("wrong status or missing result: %+v", done)
	}
	r := done.Result
	if r.TotalCents != 10000 || r.AverageCents != 2500 || r.Count != 4 {
		t.Fatalf("incorrect totals: %+v", r)
	}
	if r.Highest.Description != "Headphones" || r.Highest.AmountCents != 5000 {
		t.Fatalf("incorrect highest expense: %+v", r.Highest)
	}
	for category, want := range map[string]int64{"Food": 4000, "Transport": 1000, "Shopping": 5000} {
		if got := r.ByCategoryCents[category]; got != want {
			t.Errorf("%s: got %d, want %d", category, got, want)
		}
	}
	saved, err := store.Get(task.ID)
	if err != nil || saved.Status != StatusCompleted || saved.Result.TotalCents != 10000 {
		t.Fatalf("result not persisted: %+v, %v", saved, err)
	}
	if _, err := store.Process(task.ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected invalid state, got %v", err)
	}
}

func TestIndependentTasksAndDefensiveCopies(t *testing.T) {
	store := NewStore(ExpenseProcessor{})
	firstInput := sampleExpenses()
	a, err := store.Create(firstInput)
	if err != nil {
		t.Fatal(err)
	}
	b, err := store.Create([]Expense{{Category: "Books", Description: "Book", AmountCents: 700}})
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID {
		t.Fatal("duplicate IDs")
	}
	firstInput[0].AmountCents = 999999
	a.Expenses[0].AmountCents = 999999
	doneA, err := store.Process(a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if doneA.Result.TotalCents != 10000 {
		t.Fatal("external mutation changed stored input")
	}
	pendingB, err := store.Get(b.ID)
	if err != nil || pendingB.Status != StatusPending || pendingB.Result != nil {
		t.Fatalf("second task changed unexpectedly: %+v, %v", pendingB, err)
	}
	doneB, err := store.Process(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if doneB.Result.TotalCents != 700 {
		t.Fatalf("second result wrong: %+v", doneB.Result)
	}
	doneA.Result.ByCategoryCents["Food"] = 999
	reloadedA, err := store.Get(a.ID)
	if err != nil || reloadedA.Result.ByCategoryCents["Food"] != 4000 {
		t.Fatal("external mutation changed stored result")
	}
}
