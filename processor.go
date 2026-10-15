package main

import (
	"fmt"
	"math"
	"strings"
)

type Processor interface {
	Validate([]Expense) error
	Process([]Expense) (Result, error)
}

type ExpenseProcessor struct{}

func (ExpenseProcessor) Validate(expenses []Expense) error {
	if len(expenses) == 0 {
		return fmt.Errorf("%w: at least one expense is required", ErrInvalidInput)
	}
	var total int64
	for i, e := range expenses {
		if strings.TrimSpace(e.Category) == "" {
			return fmt.Errorf("%w: expense %d needs a category", ErrInvalidInput, i)
		}
		if strings.TrimSpace(e.Description) == "" {
			return fmt.Errorf("%w: expense %d needs a description", ErrInvalidInput, i)
		}
		if e.AmountCents <= 0 {
			return fmt.Errorf("%w: expense %d amount must be positive", ErrInvalidInput, i)
		}
		if total > math.MaxInt64-e.AmountCents {
			return fmt.Errorf("%w: total amount overflow", ErrInvalidInput)
		}
		total += e.AmountCents
	}
	return nil
}

func (p ExpenseProcessor) Process(expenses []Expense) (Result, error) {
	if err := p.Validate(expenses); err != nil {
		return Result{}, err
	}
	result := Result{
		Highest:         expenses[0],
		Count:           len(expenses),
		ByCategoryCents: make(map[string]int64),
	}
	for _, e := range expenses {
		result.TotalCents += e.AmountCents
		result.ByCategoryCents[e.Category] += e.AmountCents
		if e.AmountCents > result.Highest.AmountCents {
			result.Highest = e
		}
	}
	result.AverageCents = float64(result.TotalCents) / float64(result.Count)
	return result, nil
}
