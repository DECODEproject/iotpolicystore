ALTER TABLE policies
  ADD COLUMN uuid UUID NOT NULL;

DROP INDEX IF EXISTS policies_public_key_idx;

CREATE UNIQUE INDEX IF NOT EXISTS policies_uuid_idx
  ON policies (uuid);