-- Drop composite index on pareto_sets
DROP INDEX IF EXISTS idx_pareto_sets_user_algorithm_created;

-- Drop executions table and all its indices/constraints
DROP TABLE IF EXISTS executions;
