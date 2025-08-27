package api

import (
	"fmt"
	"net/http"

	"github.com/akarshgo/paysplit/db"
	"github.com/akarshgo/paysplit/types"
	"github.com/gofiber/fiber/v2"
)

// Wire this with your db.ExpenseStore
type ExpenseHandlers struct {
	expenses db.ExpenseStore
}

func NewExpenseHandlers(exp db.ExpenseStore) *ExpenseHandlers {
	return &ExpenseHandlers{expenses: exp}
}

// ---------- CREATE EXPENSE ----------

type createExpenseReq struct {
	PaidBy string           `json:"paid_by"`
	Note   string           `json:"note"`
	Amount int64            `json:"amount_paise"` // paise
	Split  types.SplitInput `json:"split"`
}

func (h *ExpenseHandlers) HandleCreateExpense(c *fiber.Ctx) error {
	groupID := c.Params("id")
	var req createExpenseReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if groupID == "" || req.PaidBy == "" || req.Amount <= 0 || len(req.Split.Users) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing fields"})
	}

	// 1) Normalize the split to exact amounts per user
	splits, err := normalizeSplits(req.Amount, req.Split)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 2) Build the expense row
	exp := &types.Expense{
		GroupID:     groupID,
		PaidBy:      req.PaidBy,
		AmountPaise: req.Amount,
		Currency:    "INR",
		Note:        req.Note,
		SplitKind:   req.Split.Kind,
	}

	// 3) Persist (store will create expense + insert split rows in a TX)
	id, err := h.expenses.Create(c.Context(), exp, splits)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create expense"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"id": id})
}

// ---------- LIST EXPENSES ----------

func (h *ExpenseHandlers) HandleListExpenses(c *fiber.Ctx) error {
	groupID := c.Params("id")
	out, err := h.expenses.ListByGroup(c.Context(), groupID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list expenses"})
	}
	return c.JSON(out)
}

// ---------- BALANCES & SIMPLIFY (OPTIONAL BONUS) ----------

func (h *ExpenseHandlers) HandleGroupBalances(c *fiber.Ctx) error {
	groupID := c.Params("id")
	net, err := h.expenses.Balances(c.Context(), groupID) // map[userID]int64
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to compute balances"})
	}
	return c.JSON(net)
}

type transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount_paise"`
}

func (h *ExpenseHandlers) HandleSimplifyDebts(c *fiber.Ctx) error {
	groupID := c.Params("id")
	net, err := h.expenses.Balances(c.Context(), groupID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to compute"})
	}

	var pos, neg [][2]any // [id, amt]
	for id, v := range net {
		if v > 0 {
			pos = append(pos, [2]any{id, v})
		}
		if v < 0 {
			neg = append(neg, [2]any{id, -v})
		}
	}
	i, j := 0, 0
	tx := []transfer{}
	for i < len(pos) && j < len(neg) {
		p := pos[i][1].(int64)
		n := neg[j][1].(int64)
		pay := p
		if n < p {
			pay = n
		}
		tx = append(tx, transfer{From: neg[j][0].(string), To: pos[i][0].(string), Amount: pay})
		p -= pay
		n -= pay
		pos[i][1] = p
		neg[j][1] = n
		if p == 0 {
			i++
		}
		if n == 0 {
			j++
		}
	}
	return c.JSON(tx)
}

// ---------- NORMALIZATION LOGIC ----------
// Convert client SplitInput into []ExpenseSplit with exact paise per user.
// Guarantees: len(users)>0, sum(exact) == amount, handles rounding safely.

func normalizeSplits(amount int64, in types.SplitInput) ([]types.ExpenseSplit, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be > 0")
	}
	if len(in.Users) == 0 {
		return nil, fmt.Errorf("at least one user required")
	}

	switch in.Kind {
	case types.SplitEqual:
		share := amount / int64(len(in.Users))
		rem := amount - share*int64(len(in.Users))
		out := make([]types.ExpenseSplit, 0, len(in.Users))
		for i, u := range in.Users {
			o := share
			if int64(i) < rem {
				o++
			} // distribute remainder paise fairly
			out = append(out, types.ExpenseSplit{UserID: u.UserID, Exact: types.Money(o)})
		}
		return out, nil

	case types.SplitShares:
		var sum int64
		for _, u := range in.Users {
			if u.Shares == nil || *u.Shares <= 0 {
				return nil, fmt.Errorf("shares required and must be > 0")
			}
			sum += *u.Shares
		}
		var acc int64
		out := make([]types.ExpenseSplit, len(in.Users))
		for i, u := range in.Users {
			o := (amount * *u.Shares) / sum
			out[i] = types.ExpenseSplit{UserID: u.UserID, Exact: types.Money(o)}
			acc += o
			if i == len(in.Users)-1 { // last gets remainder
				out[i].Exact += types.Money(amount - acc)
			}
		}
		return out, nil

	case types.SplitPercent:
		var sum int64
		for _, u := range in.Users {
			if u.PercentBP == nil || *u.PercentBP < 0 {
				return nil, fmt.Errorf("percent_bp required")
			}
			sum += *u.PercentBP
		}
		if sum != 10000 {
			return nil, fmt.Errorf("percent total must be exactly 10000 basis points")
		}
		var acc int64
		out := make([]types.ExpenseSplit, len(in.Users))
		for i, u := range in.Users {
			o := (amount * *u.PercentBP) / 10000
			out[i] = types.ExpenseSplit{UserID: u.UserID, Exact: types.Money(o)}
			acc += o
			if i == len(in.Users)-1 {
				out[i].Exact += types.Money(amount - acc)
			}
		}
		return out, nil

	case types.SplitExact:
		var acc int64
		for _, u := range in.Users {
			if u.Exact == nil || *u.Exact < 0 {
				return nil, fmt.Errorf("exact required and must be >= 0")
			}
			acc += *u.Exact
		}
		if acc != amount {
			return nil, fmt.Errorf("sum of exact parts must equal amount")
		}
		out := make([]types.ExpenseSplit, len(in.Users))
		for i, u := range in.Users {
			out[i] = types.ExpenseSplit{UserID: u.UserID, Exact: types.Money(*u.Exact)}
		}
		return out, nil
	}

	return nil, fmt.Errorf("unknown split kind: %s", in.Kind)
}
