package types

import "time"

type Money int64 // paise

type User struct {
	ID        string
	Name      string
	Email     *string
	Phone     *string
	UPI       *string // e.g. akarsh@okicici
	CreatedAt time.Time
}

type Group struct {
	ID        string
	Name      string
	CreatedBy string // UserID
	CreatedAt time.Time
}

type SplitPart struct {
	UserID string
	Share  int   // numerator for ratio splits
	Exact  Money // optional for exact splits
}

type Settlement struct {
	ID         string
	FromUserID string
	ToUserID   string
	Amount     Money
	Method     string  // "UPI","Cash","Bank"
	Ref        *string // UPI txn id, etc.
	CreatedAt  time.Time
}

type Reminder struct {
	ID        string
	DebtKey   string // composite: debtor->creditor
	NextAt    time.Time
	Frequency string // "DAILY","WEEKLY","CUSTOM"
	Channel   string // "PUSH","SMS","WHATSAPP","EMAIL"
	Active    bool
	CreatedBy string // who set it
	CreatedAt time.Time
}
