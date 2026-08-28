package model

import "slices"

// Guest represents a wedding guest.
type Guest struct {
	ID              string  `gorm:"primaryKey;notnull"    json:"id,omitempty"`
	FirstName       string  `gorm:"uniqueIndex:idx_name"  json:"first_name,omitempty"`
	LastName        string  `gorm:"uniqueIndex:idx_name"  json:"last_name,omitempty"`
	SpouseFirstName *string `gorm:"uniqueIndex:idx_sname" json:"spouse_first_name,omitempty"`
	SpouseLastName  *string `gorm:"uniqueIndex:idx_sname" json:"spouse_last_name,omitempty"`
	IsFamily        bool    `json:"is_family,omitempty"`
	Children        int     `json:"childeren,omitempty"`
	// Events is the comma-separated set of ceremonies this guest is invited
	// to. Empty means all of them, so a guest list that predates tiering, or
	// a wedding that does not tier at all, needs no migration.
	Events string  `json:"events,omitempty"`
	Answer *Answer `gorm:"foreignKey:GuestID" json:"answer,omitempty"`
}

// InvitedEvents resolves the guest's events, expanding the empty set.
func (g Guest) InvitedEvents() []Event {
	events, err := ParseEvents(g.Events)
	if err != nil {
		// A stored value can only be invalid if it was written outside the
		// service; fall back to the safe reading rather than dropping them.
		return AllEvents()
	}

	return events
}

// InvitedTo reports whether the guest is invited to a given ceremony.
func (g Guest) InvitedTo(event Event) bool {
	return slices.Contains(g.InvitedEvents(), event)
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
