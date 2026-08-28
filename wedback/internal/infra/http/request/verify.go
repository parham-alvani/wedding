package request

// Verify carries the name a visitor typed to open an invitation.
type Verify struct {
	Name string `json:"name,omitempty"`
}
