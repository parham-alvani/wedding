package model

import "time"

// Guest represents a wedding guest that may also has partner.
type Guest struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Partner   *Guest    `json:"partner"`
	Answer    *Answer   `json:"answer"`
}
