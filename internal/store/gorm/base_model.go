package gorm

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains the common columns for all tables.
type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
