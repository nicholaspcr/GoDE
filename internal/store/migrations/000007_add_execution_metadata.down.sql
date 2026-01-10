-- Remove algorithm, variant, and problem columns from executions table
ALTER TABLE executions DROP COLUMN algorithm;
ALTER TABLE executions DROP COLUMN variant;
ALTER TABLE executions DROP COLUMN problem;
