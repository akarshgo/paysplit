package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/akarshgo/paysplit/types"
	"github.com/google/uuid"
)

type ExpenseStore interface {
	Create(ctx context.Context, e *types.Expense, splits []types.ExpenseSplit) (string, error)
	ListByGroup(ctx context.Context, groupID string) ([]*types.Expense, error)
	Balances(ctx context.Context, groupID string) (map[string]int64, error)
}

type PostgresExpenseStore struct {
	db *sql.DB
}

func NewPostgresExpenseStore(db *sql.DB) *PostgresExpenseStore {
	return &PostgresExpenseStore{db: db}
}

// Create inserts expense + splits in one transaction
func (s *PostgresExpenseStore) Create(ctx context.Context, e *types.Expense, splits []types.ExpenseSplit) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	id := uuid.New().String()
	now := time.Now()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO expenses (id, group_id, paid_by, amount_paise, currency, note, split_kind, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, id, e.GroupID, e.PaidBy, e.AmountPaise, e.Currency, e.Note, string(e.SplitKind), now)
	if err != nil {
		return "", err
	}

	for _, sp := range splits {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO expense_splits (id, expense_id, user_id, exact)
			VALUES ($1,$2,$3,$4)
		`, uuid.New().String(), id, sp.UserID, sp.Exact)
		if err != nil {
			return "", err
		}
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	e.ID = id
	e.CreatedAt = now
	return id, nil
}

// ListByGroup fetches all expenses for a group
func (s *PostgresExpenseStore) ListByGroup(ctx context.Context, groupID string) ([]*types.Expense, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, group_id, paid_by, amount_paise, currency, note, split_kind, created_at
		FROM expenses
		WHERE group_id = $1
		ORDER BY created_at DESC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*types.Expense
	for rows.Next() {
		var e types.Expense
		if err := rows.Scan(&e.ID, &e.GroupID, &e.PaidBy, &e.AmountPaise, &e.Currency, &e.Note, &e.SplitKind, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

// Balances computes net balance per user in group
func (s *PostgresExpenseStore) Balances(ctx context.Context, groupID string) (map[string]int64, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT e.paid_by, s.user_id, s.exact
		FROM expenses e
		JOIN expense_splits s ON s.expense_id = e.id
		WHERE e.group_id = $1
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	net := map[string]int64{}
	for rows.Next() {
		var paidBy, userID string
		var exact int64
		if err := rows.Scan(&paidBy, &userID, &exact); err != nil {
			return nil, err
		}
		net[paidBy] += exact // payer gets credit
		net[userID] -= exact // participant owes
	}
	return net, rows.Err()
}
