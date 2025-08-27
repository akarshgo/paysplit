package db

import (
	"context"
	"database/sql"

	"github.com/akarshgo/paysplit/types"
)

type ExpenseStore interface {
	Create(ctx context.Context, e *types.Expense, splits []types.ExpenseSplit) (string, error)
	ListByGroup(ctx context.Context, groupID string) (*[]types.Expense, error)
	Balances(ctx context.Context, groupID string) (map[string]int64, error) //net paise par user
}

type PostGresExpsenseStore struct {
	db *sql.DB
}

func NewPostGresExpenseStore(db *sql.DB) *PostGresExpsenseStore {
	return &PostGresExpsenseStore{
		db: db,
	}
}

func (p *PostGresExpsenseStore) Create(ctx context.Context, e *types.Expense, splits []types.ExpenseSplit) (string, error) {
	return "", nil
}

func (p *PostGresExpsenseStore) ListByGroup(ctx context.Context, groupID string) (*[]types.Expense, error) {
	return nil, nil
}

func (p *PostGresExpsenseStore) Balances(ctx context.Context, groupID string) (map[string]int64, error) {
	return nil, nil

}
