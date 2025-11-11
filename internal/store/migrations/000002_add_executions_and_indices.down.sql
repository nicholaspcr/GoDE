-- Rollback JSONB conversion (PostgreSQL only)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_catalog.pg_type WHERE typname = 'jsonb') THEN
        -- Drop GIN indices
        DROP INDEX IF EXISTS idx_vectors_objectives_gin;
        DROP INDEX IF EXISTS idx_vectors_elements_gin;

        -- Convert JSONB back to TEXT
        ALTER TABLE vectors
            ALTER COLUMN objectives TYPE TEXT USING objectives::TEXT;

        ALTER TABLE vectors
            ALTER COLUMN elements TYPE TEXT USING elements::TEXT;
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        -- Ignore errors for non-PostgreSQL databases
        NULL;
END$$;

-- Drop composite index on pareto_sets
DROP INDEX IF EXISTS idx_pareto_sets_user_algorithm_created;

-- Drop executions table and all its indices/constraints
DROP TABLE IF EXISTS executions CASCADE;
