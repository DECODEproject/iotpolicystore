CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS policies (
  id SERIAL PRIMARY KEY,
  public_key TEXT NOT NULL,
  label TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  operations JSONB,
  token BYTEA NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS policies_public_key_idx
  ON policies (public_key);