-- Initial schema for GoDE

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Pareto sets table
CREATE TABLE IF NOT EXISTS pareto_sets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    algorithm VARCHAR(100) NOT NULL,
    problem VARCHAR(100) NOT NULL,
    variant VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pareto_sets_user_id ON pareto_sets(user_id);
CREATE INDEX IF NOT EXISTS idx_pareto_sets_algorithm ON pareto_sets(algorithm);
CREATE INDEX IF NOT EXISTS idx_pareto_sets_created_at ON pareto_sets(created_at);

-- Vectors table (Pareto solutions)
CREATE TABLE IF NOT EXISTS vectors (
    id SERIAL PRIMARY KEY,
    pareto_set_id INTEGER NOT NULL REFERENCES pareto_sets(id) ON DELETE CASCADE,
    elements TEXT NOT NULL,
    objectives TEXT NOT NULL,
    crowding_distance DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_vectors_pareto_set_id ON vectors(pareto_set_id);
