package main

import "fmt"

func main() {
	store := NewStore(ExpenseProcessor{})
	inputs := [][]Expense{
		{
			{Category: "Food", Description: "Lunch", AmountCents: 1500},
			{Category: "Transport", Description: "Taxi", AmountCents: 1000},
			{Category: "Food", Description: "Dinner", AmountCents: 2500},
			{Category: "Shopping", Description: "Headphones", AmountCents: 5000},
		},
		{
			{Category: "Books", Description: "Go book", AmountCents: 3000},
			{Category: "Food", Description: "Snack", AmountCents: 500},
		},
	}
	for _, expenses := range inputs {
		task, err := store.Create(expenses)
		if err != nil {
			fmt.Println("create error:", err)
			continue
		}
		fmt.Printf("Created task %d (%s)\n", task.ID, task.Status)
		task, err = store.Process(task.ID)
		if err != nil {
			fmt.Println("processing error:", err)
			continue
		}
		fmt.Printf("Task %d (%s): total $%.2f, average $%.2f, highest %s $%.2f\n",
			task.ID, task.Status,
			float64(task.Result.TotalCents)/100,
			task.Result.AverageCents/100,
			task.Result.Highest.Description,
			float64(task.Result.Highest.AmountCents)/100)
		fmt.Println("By category (cents):", task.Result.ByCategoryCents)
	}
}
