package db

import (
	"context"
	"database/sql"

	"github.com/akarshgo/paysplit/types"
)

type GroupStore interface {
	Create(ctx context.Context, g *types.Group) (string, error)
	ListByUser(ctx context.Context, userID string) ([]*types.Group, error)
	AddMember(ctx context.Context, groupID, userID string) error
}

type PostGresGroupStore struct {
	db *sql.DB
}

func NewPostGresGroupStore(db *sql.DB) *PostGresGroupStore {
	return &PostGresGroupStore{
		db: db,
	}
}

func (p *PostGresGroupStore) Create(ctx context.Context, g *types.Group) (string, error) {
	return "", nil
}

func (p *PostGresGroupStore) ListByUser(ctx context.Context, userID string) ([]*types.Group, error) {
	return nil, nil
}

func (p *PostGresGroupStore) AddMember(ctx context.Context, groupID, userID string) error {
	return nil
}
