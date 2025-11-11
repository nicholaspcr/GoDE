-- Add executions table for async DE execution tracking
CREATE TABLE IF NOT EXISTS executions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    config_json TEXT NOT NULL,
    pareto_id BIGINT,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Individual indices for executions table
CREATE INDEX IF NOT EXISTS idx_executions_user_id ON executions(user_id);
CREATE INDEX IF NOT EXISTS idx_executions_status ON executions(status);
CREATE INDEX IF NOT EXISTS idx_executions_created_at ON executions(created_at);
CREATE INDEX IF NOT EXISTS idx_executions_pareto_id ON executions(pareto_id);

-- Composite index for optimized ListExecutions query (user + status + created_at)
CREATE INDEX IF NOT EXISTS idx_executions_user_status_created ON executions(user_id, status, created_at DESC);

-- Foreign key constraint to pareto_sets
ALTER TABLE executions ADD CONSTRAINT fk_executions_pareto
    FOREIGN KEY (pareto_id) REFERENCES pareto_sets(id) ON DELETE SET NULL;

-- Add composite index on pareto_sets for common query patterns
CREATE INDEX IF NOT EXISTS idx_pareto_sets_user_algorithm_created ON pareto_sets(user_id, algorithm, created_at DESC);

-- Optimize vectors table with JSONB columns (PostgreSQL only)
-- For SQLite, this will be ignored
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_catalog.pg_type WHERE typname = 'jsonb') THEN
        -- Convert TEXT columns to JSONB for better querying
        ALTER TABLE vectors
            ALTER COLUMN elements TYPE JSONB USING elements::JSONB;

        ALTER TABLE vectors
            ALTER COLUMN objectives TYPE JSONB USING objectives::JSONB;

        -- Add GIN indices for JSONB columns
        CREATE INDEX IF NOT EXISTS idx_vectors_elements_gin ON vectors USING GIN (elements);
        CREATE INDEX IF NOT EXISTS idx_vectors_objectives_gin ON vectors USING GIN (objectives);
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        -- Ignore errors for non-PostgreSQL databases
        NULL;
END$$;
