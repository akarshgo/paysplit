package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/akarshgo/paysplit/types"
	"github.com/google/uuid"
)

type GroupStore interface {
	Create(ctx context.Context, g *types.Group) (string, error)
	ListByUser(ctx context.Context, userID string) ([]*types.Group, error)
	AddMember(ctx context.Context, groupID, userID string) error
}

type PostgresGroupStore struct {
	db *sql.DB
}

func NewPostgresGroupStore(db *sql.DB) *PostgresGroupStore {
	return &PostgresGroupStore{db: db}
}

func (p *PostgresGroupStore) Create(ctx context.Context, g *types.Group) (string, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := p.db.ExecContext(ctx, `
		INSERT INTO groups (id, name, created_by, created_at)
		VALUES ($1, $2, $3, $4)
	`, id, g.Name, g.CreatedBy, now)
	if err != nil {
		return "", err
	}

	// also insert creator as member (admin)
	_, _ = p.db.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id, role, added_at)
		VALUES ($1, $2, 'admin', $3)
		ON CONFLICT DO NOTHING
	`, id, g.CreatedBy, now)

	g.ID, g.CreatedAt = id, now
	return id, nil
}

func (p *PostgresGroupStore) ListByUser(ctx context.Context, userID string) ([]*types.Group, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT g.id, g.name, g.created_by, g.created_at
		FROM groups g
		JOIN group_members m ON m.group_id = g.id
		WHERE m.user_id = $1
		ORDER BY g.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*types.Group
	for rows.Next() {
		var g types.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, rows.Err()
}

func (p *PostgresGroupStore) AddMember(ctx context.Context, groupID, userID string) error {
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id, role, added_at)
		VALUES ($1, $2, 'member', $3)
		ON CONFLICT DO NOTHING
	`, groupID, userID, time.Now())
	return err
}
