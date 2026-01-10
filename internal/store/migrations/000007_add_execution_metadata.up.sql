-- Add algorithm, variant, and problem columns to executions table
ALTER TABLE executions ADD COLUMN algorithm VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE executions ADD COLUMN variant VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE executions ADD COLUMN problem VARCHAR(50) NOT NULL DEFAULT '';
