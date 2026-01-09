-- Remove deleted_at columns

DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;

DROP INDEX IF EXISTS idx_pareto_sets_deleted_at;
ALTER TABLE pareto_sets DROP COLUMN deleted_at;

DROP INDEX IF EXISTS idx_vectors_deleted_at;
ALTER TABLE vectors DROP COLUMN deleted_at;
