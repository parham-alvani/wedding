package model

// Answer shows the answer people said about coming to our wedding or not.
type Answer struct {
	ID      int64 `gorm:"primaryKey;autoIncrement;not null"`
	Coming  bool
	PlusOne bool
	GuestID string `gorm:"uniqueIndex;not null"`
}
