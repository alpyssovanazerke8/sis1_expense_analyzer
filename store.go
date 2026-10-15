package main

import "fmt"

// Store is synchronous and in-memory. It is not concurrency-safe.
type Store struct {
	nextID    int
	tasks     map[int]*Task
	processor Processor
}

func NewStore(processor Processor) *Store {
	return &Store{nextID: 1, tasks: make(map[int]*Task), processor: processor}
}

func (s *Store) Create(expenses []Expense) (Task, error) {
	if err := s.processor.Validate(expenses); err != nil {
		return Task{}, err
	}
	task := &Task{
		ID:       s.nextID,
		Expenses: append([]Expense(nil), expenses...),
		Status:   StatusPending,
	}
	s.tasks[task.ID] = task
	s.nextID++
	return task.Clone(), nil
}

func (s *Store) Get(id int) (Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	return task.Clone(), nil
}

func (s *Store) Process(id int) (Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	if task.Status != StatusPending {
		return Task{}, fmt.Errorf("%w: cannot process task in %s state", ErrInvalidState, task.Status)
	}
	task.Status = StatusProcessing
	result, err := s.processor.Process(task.Expenses)
	if err != nil {
		// Restore pending state so a failed attempt does not appear completed.
		task.Status = StatusPending
		return Task{}, err
	}
	task.Result = &result
	task.Status = StatusCompleted
	return task.Clone(), nil
}
