-- Add updated_at column to vectors table (required by gorm.Model)
ALTER TABLE vectors ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
