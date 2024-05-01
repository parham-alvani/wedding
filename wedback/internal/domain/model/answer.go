package model

// Answer shows the answer people said about coming to our wedding or not.
type Answer struct {
	ID      int64  `gorm:"primaryKey;autoIncrement;not null" json:"-"`
	Coming  bool   `json:"coming"`
	PlusOne bool   `json:"plus_one"`
	GuestID string `gorm:"uniqueIndex;not null"              json:"-"`
}
