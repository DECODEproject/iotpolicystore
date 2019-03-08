ALTER TABLE policies
  ADD COLUMN authorizable_attribute_id TEXT NOT NULL,
  ADD COLUMN credential_issuer_endpoint_url TEXT NOT NULL;
