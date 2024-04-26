package model

import (
	"gorm.io/gorm"
)

// Guest represents a wedding guest that may also has partner.
type Guest struct {
	gorm.Model
	ID     string  `gorm:"primaryKey" json:"id"`
	Name   string  `json:"name"`
	Answer *Answer `gorm:"foreignKey:ID" json:"answer"`
}
