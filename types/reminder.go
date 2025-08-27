package types

import "time"

type Reminder struct {
	ID         string    `json:"id"`
	DebtKey    string    `json:"debt_key"` // e.g., groupID:userID pair
	TargetUser string    `json:"target_user"`
	NextAt     time.Time `json:"next_at"`
	Frequency  string    `json:"frequency"` // "daily","weekly"
	Channel    string    `json:"channel"`   // "push","sms","email"
	Active     bool      `json:"active"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}
