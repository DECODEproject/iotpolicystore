CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS policies (
  id SERIAL PRIMARY KEY,
  public_key TEXT NOT NULL CHECK (public_key <> ''),
  label TEXT NOT NULL CHECK (label <> ''),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  operations JSONB,
  token BYTEA NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS policies_public_key_idx
  ON policies (public_key);