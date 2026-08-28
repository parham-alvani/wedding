package wedding

// Config holds wedding-specific configuration.
type Config struct {
	HusbandName string `json:"husband_name,omitempty" koanf:"husband_name"`
	WifeName    string `json:"wife_name,omitempty"    koanf:"wife_name"`
	BaseURL     string `json:"base_url,omitempty"     koanf:"base_url"`

	// RSVPDeadline is an RFC 3339 timestamp after which guests can no longer
	// answer or change their answer. Empty means the RSVP never closes.
	RSVPDeadline string `json:"rsvp_deadline,omitempty" koanf:"rsvp_deadline"`

	// InviteTemplate is the message `wedback invite` renders for each guest.
	// It is a Go text/template over the fields documented on that command.
	InviteTemplate string `json:"invite_template,omitempty" koanf:"invite_template"`
}
