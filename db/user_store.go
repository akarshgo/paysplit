package db

import (
	"context"
	"database/sql"

	"github.com/akarshgo/paysplit/types"
)

// store/user_store.go
type UserFilter struct {
	Email *string
	Phone *string
	Query *string // search by name/email/phone
}

type UserStore interface {
	Create(ctx context.Context, u *types.User) error
	GetByID(ctx context.Context, id string) (*types.User, error)
	Find(ctx context.Context, f UserFilter, limit, offset int) ([]*types.User, error)
	Update(ctx context.Context, u *types.User) error
	UpsertUPI(ctx context.Context, userID, vpa string) error
	Delete(ctx context.Context, id string) error
	// auth helpers
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	GetByPhone(ctx context.Context, phone string) (*types.User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}
