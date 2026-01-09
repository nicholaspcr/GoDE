-- Add deleted_at column for soft deletes (required by gorm.Model)

ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

ALTER TABLE pareto_sets ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX IF NOT EXISTS idx_pareto_sets_deleted_at ON pareto_sets(deleted_at);

ALTER TABLE vectors ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX IF NOT EXISTS idx_vectors_deleted_at ON vectors(deleted_at);
