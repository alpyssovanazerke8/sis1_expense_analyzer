package main

import "errors"

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
)

var (
	ErrInvalidInput = errors.New("invalid expense input")
	ErrNotFound     = errors.New("task not found")
	ErrInvalidState = errors.New("invalid task state")
)

type Expense struct {
	Category    string
	Description string
	AmountCents int64
}

type Result struct {
	TotalCents      int64
	AverageCents    float64
	Highest         Expense
	ByCategoryCents map[string]int64
	Count           int
}

type Task struct {
	ID       int
	Expenses []Expense
	Status   Status
	Result   *Result
}

// Clone ensures callers cannot modify stored tasks through returned slices or maps.
func (t Task) Clone() Task {
	t.Expenses = append([]Expense(nil), t.Expenses...)
	if t.Result != nil {
		result := *t.Result
		result.ByCategoryCents = make(map[string]int64, len(t.Result.ByCategoryCents))
		for k, v := range t.Result.ByCategoryCents {
			result.ByCategoryCents[k] = v
		}
		t.Result = &result
	}
	return t
}
