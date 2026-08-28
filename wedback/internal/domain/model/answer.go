package model

// Answer shows the answer people said about coming to our wedding or not.
type Answer struct {
	ID      int64 `gorm:"primaryKey;autoIncrement;not null" json:"-"`
	Coming  bool  `json:"coming"`
	PlusOne bool  `json:"plus_one"`
	// Dietary is anything the kitchen needs to know: allergies, vegetarian,
	// and so on. Free text, because real requirements never fit a checkbox.
	Dietary string `json:"dietary,omitempty"`
	// Song is what this guest would like to hear on the night.
	Song    string `json:"song,omitempty"`
	GuestID string `gorm:"uniqueIndex;not null" json:"-"`
}
