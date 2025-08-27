package types

type SplitKind string

const (
	SplitEqual   SplitKind = "equal"
	SplitShares  SplitKind = "shares"  // weighted shares (1,1,2)
	SplitPercent SplitKind = "percent" // basis points (10000 = 100.00%)
	SplitExact   SplitKind = "exact"   // client provides exact paise per user
)

type SplitInputUser struct {
	UserID    string `json:"user_id"`
	Shares    *int64 `json:"shares,omitempty"`     // for shares
	PercentBP *int64 `json:"percent_bp,omitempty"` // for percent
	Exact     *int64 `json:"exact,omitempty"`      // for exact
}

type SplitInput struct {
	Kind  SplitKind        `json:"kind"`
	Users []SplitInputUser `json:"users"`
}
