package model

// Guest represents a wedding guest.
type Guest struct {
	ID              string  `gorm:"primaryKey;notnull"    json:"id,omitempty"`
	FirstName       string  `gorm:"uniqueIndex:idx_name"  json:"first_name,omitempty"`
	LastName        string  `gorm:"uniqueIndex:idx_name"  json:"last_name,omitempty"`
	SpouseFirstName *string `gorm:"uniqueIndex:idx_sname" json:"spouse_first_name,omitempty"`
	SpouseLastName  *string `gorm:"uniqueIndex:idx_sname" json:"spouse_last_name,omitempty"`
	IsFamily        bool    `json:"is_family,omitempty"`
	Children        int     `json:"childeren,omitempty"`
	Answer          *Answer `gorm:"foreignKey:GuestID"    json:"answer,omitempty"`
}

func (g Guest) Coming() bool {
	if g.IsFamily {
		return true
	}

	if g.Answer != nil {
		return g.Answer.Coming
	}

	return false
}

func (g Guest) PlusOne() bool {
	if g.IsFamily && g.SpouseFirstName != nil {
		return true
	}

	if g.Answer != nil {
		return g.Answer.PlusOne
	}

	return false
}
