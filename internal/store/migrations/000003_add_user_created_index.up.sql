-- Optimize: SELECT * FROM executions WHERE user_id = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_executions_user_created
ON executions(user_id, created_at DESC);

-- Also optimize pareto queries
CREATE INDEX IF NOT EXISTS idx_pareto_sets_user_created
ON pareto_sets(user_id, created_at DESC);
