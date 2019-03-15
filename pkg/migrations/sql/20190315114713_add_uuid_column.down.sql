ALTER TABLE policies
  DROP COLUMN uuid;

CREATE UNIQUE INDEX IF NOT EXISTS policies_public_key_idx
  ON policies (public_key);