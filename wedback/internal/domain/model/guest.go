package model

import "time"

// Guest represents a wedding guest that may also has partner.
type Guest struct {
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`
	ID        string    `gorm:"primaryKey" yaml:"id"`
	Name      string    `yaml:"name"`
	Email     string    `yaml:"email"`
	Partner   *Guest    `yaml:"partner"`
	Answer    *Answer   `yaml:"answer"`
}
