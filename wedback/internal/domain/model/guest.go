package model

// Guest represents a wedding guest.
type Guest struct {
	ID              string  `gorm:"primaryKey;notnull"    json:"id"`
	FirstName       string  `gorm:"uniqueIndex:idx_name"  json:"first_name"`
	LastName        string  `gorm:"uniqueIndex:idx_name"  json:"last_name"`
	SpouseFirstName *string `gorm:"uniqueIndex:idx_sname" json:"spouse_first_name,omitempty"`
	SpouseLastName  *string `gorm:"uniqueIndex:idx_sname" json:"spouse_last_name,omitempty"`
	Answer          *Answer `gorm:"foreignKey:GuestID"    json:"answer"`
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
