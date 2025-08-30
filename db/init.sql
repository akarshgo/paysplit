CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIUQE,
    phone TEXT UNIUQE,
    upi_vpa TEXT,
    creted_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS groups(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    creted_by UUID REFERENCES users(id),
    creted_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role TEXT DEFAULT `member`
    added_at TIMESTAMPZ NOT NULL DEFAULT now()
    PRIMARY KEY (group_id, user_id)
)


CREATE TABLE IF NOT EXISTS expenses (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id      UUID REFERENCES groups(id) ON DELETE CASCADE,
  paid_by       UUID REFERENCES users(id),
  amount_paise  BIGINT NOT NULL,
  currency      TEXT   NOT NULL DEFAULT 'INR',
  note          TEXT,
  split_kind    TEXT   NOT NULL,   -- equal|shares|percent|exact
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS expense_splits (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  expense_id UUID REFERENCES expenses(id) ON DELETE CASCADE,
  user_id    UUID REFERENCES users(id),
  exact      BIGINT NOT NULL
);

-- helpful indexes
CREATE INDEX IF NOT EXISTS idx_expenses_group ON expenses(group_id);
CREATE INDEX IF NOT EXISTS idx_splits_expense ON expense_splits(expense_id);
CREATE INDEX IF NOT EXISTS idx_splits_user    ON expense_splits(user_id);
