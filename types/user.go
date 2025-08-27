package types

import "time"

// Nullable fields are pointers so we can distinguish "unset" vs "".
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	UPI       *string   `json:"upi_vpa,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
