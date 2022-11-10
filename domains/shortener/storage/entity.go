package storage

import (
	"gorm.io/gorm"
)

// Record sample record for
type Record struct {
	ID        uint64 `gorm:"primaryKey"`
	LongText  string `json:"long"`
	ShortText string `json:"short" gorm:"uniqueIndex;size:10"`
	gorm.Model
}
