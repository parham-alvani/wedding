package model

import "gorm.io/gorm"

// Answer shows the answer people said about coming to our wedding or not.
type Answer struct {
	gorm.Model

	Coming  bool
	PlusOne bool
}
