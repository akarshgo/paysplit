package types

import "time"

// Money is paise (₹1.00 => 100)

type Expense struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	PaidBy      string    `json:"paid_by"`
	AmountPaise int64     `json:"amount_paise"`
	Currency    string    `json:"currency"` // "INR"
	Note        string    `json:"note"`
	SplitKind   SplitKind `json:"split_kind"`
	CreatedAt   time.Time `json:"created_at"`
}

// What we insert into expense_splits (already normalized to exact paise)
type ExpenseSplit struct {
	ID        string `json:"id,omitempty"`
	ExpenseID string `json:"expense_id,omitempty"`
	UserID    string `json:"user_id"`
	Exact     Money  `json:"exact"` // paise each user owes for this expense
}
