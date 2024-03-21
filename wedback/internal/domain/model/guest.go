package model

import "time"

// Guest represents a wedding guest that may also has partner.
type Guest struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string
	Partner   *Guest
	Answer    *Answer
}
