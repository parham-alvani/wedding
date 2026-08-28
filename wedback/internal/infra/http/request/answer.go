package request

type Answer struct {
	PlusOne bool   `json:"plus_one,omitempty"`
	Coming  bool   `json:"coming,omitempty"`
	Dietary string `json:"dietary,omitempty"`
	Song    string `json:"song,omitempty"`
}
