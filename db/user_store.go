package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akarshgo/paysplit/types"
	"github.com/google/uuid"
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

func (s *PostgresUserStore) Create(ctx context.Context, u *types.User) error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name required")
	}
	id := uuid.New().String()
	now := time.Now()

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, phone, upi_vpa, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, id, u.Name, nullStr(u.Email), nullStr(u.Phone), nullStr(u.UPI), now)
	if err != nil {
		return err
	}
	u.ID = id
	u.CreatedAt = now
	return nil
}

func (s *PostgresUserStore) GetByID(ctx context.Context, id string) (*types.User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, phone, upi_vpa, created_at
		FROM users WHERE id = $1
	`, id)
	return scanUser(row)
}

func (s *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, phone, upi_vpa, created_at
		FROM users WHERE email = $1
	`, email)
	return scanUser(row)
}

func (s *PostgresUserStore) GetByPhone(ctx context.Context, phone string) (*types.User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, phone, upi_vpa, created_at
		FROM users WHERE phone = $1
	`, phone)
	return scanUser(row)
}

func (s *PostgresUserStore) Find(ctx context.Context, f UserFilter, limit, offset int) ([]*types.User, error) {
	var (
		where []string
		args  []any
		i     = 1
	)
	if f.Email != nil {
		where = append(where, fmt.Sprintf("email = $%d", i))
		args = append(args, *f.Email)
		i++
	}
	if f.Phone != nil {
		where = append(where, fmt.Sprintf("phone = $%d", i))
		args = append(args, *f.Phone)
		i++
	}
	if f.Query != nil && strings.TrimSpace(*f.Query) != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)", i, i, i))
		args = append(args, "%"+strings.TrimSpace(*f.Query)+"%")
		i++
	}
	q := `
		SELECT id, name, email, phone, upi_vpa, created_at
		FROM users`
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	q += " ORDER BY created_at DESC"
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*types.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *PostgresUserStore) Update(ctx context.Context, u *types.User) error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("id required")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name required")
	}
	_, err := s.db.ExecContext(ctx, `
		UPDATE users
		SET name = $2, email = $3, phone = $4, upi_vpa = $5
		WHERE id = $1
	`, u.ID, u.Name, nullStr(u.Email), nullStr(u.Phone), nullStr(u.UPI))
	return err
}

func (s *PostgresUserStore) UpsertUPI(ctx context.Context, userID, vpa string) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE users SET upi_vpa = $2 WHERE id = $1
	`, userID, vpa)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresUserStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

// --- helpers ---

func scanUser(scanner interface{ Scan(dest ...any) error }) (*types.User, error) {
	var (
		u         types.User
		emailNS   sql.NullString
		phoneNS   sql.NullString
		upiNS     sql.NullString
		createdAt time.Time
	)
	if err := scanner.Scan(&u.ID, &u.Name, &emailNS, &phoneNS, &upiNS, &createdAt); err != nil {
		return nil, err
	}
	if emailNS.Valid {
		u.Email = &emailNS.String
	}
	if phoneNS.Valid {
		u.Phone = &phoneNS.String
	}
	if upiNS.Valid {
		u.UPI = &upiNS.String
	}
	u.CreatedAt = createdAt
	return &u, nil
}

func nullStr(s *string) any {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
