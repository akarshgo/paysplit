package types

import "time"

type Settlement struct {
	ID        string    `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Amount    Money     `json:"amount"`
	Method    string    `json:"method"` // "upi", "cash", ...
	Ref       string    `json:"ref,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
