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

func (g Guest) Coming() bool {
	if g.Answer != nil {
		return g.Answer.Coming
	}
	return false
}

func (g Guest) PlusOne() bool {
	if g.Answer != nil {
		return g.Answer.PlusOne
	}
	return false
}
