package model

// Guest represents a wedding guest.
type Guest struct {
	ID     string  `gorm:"primaryKey;notnull" json:"id"`
	Name   string  `gorm:"uniqueIndex"        json:"name"`
	Spouse string  `json:"spouse"`
	Answer *Answer `gorm:"foreignKey:GuestID" json:"answer"`
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
