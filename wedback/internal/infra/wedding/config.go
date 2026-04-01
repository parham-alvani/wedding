package wedding

// Config holds wedding-specific configuration.
type Config struct {
	HusbandName string `json:"husband_name,omitempty" koanf:"husband_name"`
	WifeName    string `json:"wife_name,omitempty"    koanf:"wife_name"`
	BaseURL     string `json:"base_url,omitempty"     koanf:"base_url"`
}
