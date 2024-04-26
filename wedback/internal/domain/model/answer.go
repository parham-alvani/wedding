package model

// Answer shows the answer people said about coming to our wedding or not.
type Answer struct {
	ID      int64  `gorm:"primaryKey;autoIncrement;not null" json:"-"`
	Coming  bool   `json:"coming,omitempty"`
	PlusOne bool   `json:"plus_one,omitempty"`
	GuestID string `gorm:"uniqueIndex;not null" json:"-"`
}
