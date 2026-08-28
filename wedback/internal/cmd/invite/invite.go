// Package invite turns the guest list into per-guest invitation messages and
// QR codes. Every guest has their own RSVP link, and delivering a hundred of
// them by hand is the most tedious part of running the site.
package invite

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	qrcode "github.com/skip2/go-qrcode"
)

// DefaultTemplate is used when wedding.invite_template is not set.
const DefaultTemplate = `Dear {{ .Names }},

{{ .Wife }} and {{ .Husband }} would love to have you at their wedding.
Please let us know whether you can make it: {{ .Link }}`

// qrSize is the width of a generated QR PNG, in pixels. Large enough to
// survive being printed on a card.
const qrSize = 512

var ErrNoGuests = errors.New("no guests matched")

// Invitation is the data a template is rendered against.
type Invitation struct {
	ID              string
	FirstName       string
	LastName        string
	SpouseFirstName string
	SpouseLastName  string
	// Names is "Ali" or "Ali & Maryam", whichever fits the guest.
	Names    string
	IsFamily bool
	Children int
	Link     string
	// Husband and Wife are the couple's names, from configuration.
	Husband string
	Wife    string
	// Answered reports whether this guest has already replied.
	Answered bool
	Coming   bool
}

// Rendered pairs a guest with their finished message.
type Rendered struct {
	Invitation Invitation
	Message    string
}

// Link builds a guest's personal invitation URL.
func Link(baseURL, id string) string {
	return strings.TrimSuffix(baseURL, "/") + "/guests/" + id
}

// NewInvitation projects a guest into template data.
func NewInvitation(guest model.Guest, baseURL, husband, wife string) Invitation {
	names := guest.FirstName
	if guest.SpouseFirstName != nil && *guest.SpouseFirstName != "" {
		names = guest.FirstName + " & " + *guest.SpouseFirstName
	}

	spouseFirst, spouseLast := "", ""
	if guest.SpouseFirstName != nil {
		spouseFirst = *guest.SpouseFirstName
	}

	if guest.SpouseLastName != nil {
		spouseLast = *guest.SpouseLastName
	}

	return Invitation{
		ID:              guest.ID,
		FirstName:       guest.FirstName,
		LastName:        guest.LastName,
		SpouseFirstName: spouseFirst,
		SpouseLastName:  spouseLast,
		Names:           names,
		IsFamily:        guest.IsFamily,
		Children:        guest.Children,
		Link:            Link(baseURL, guest.ID),
		Husband:         husband,
		Wife:            wife,
		Answered:        guest.Answer != nil,
		Coming:          guest.Coming(),
	}
}

// Waiting keeps only the guests who still owe an answer. Family guests are
// counted as attending without replying, so they are never waiting.
func Waiting(guests []model.Guest) []model.Guest {
	out := make([]model.Guest, 0, len(guests))

	for _, guest := range guests {
		if guest.Answer == nil && !guest.IsFamily {
			out = append(out, guest)
		}
	}

	return out
}

// Render renders one message per guest. An empty tmpl uses DefaultTemplate.
func Render(tmpl string, guests []model.Guest, baseURL, husband, wife string) ([]Rendered, error) {
	if strings.TrimSpace(tmpl) == "" {
		tmpl = DefaultTemplate
	}

	parsed, err := template.New("invite").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("invite template is invalid: %w", err)
	}

	if len(guests) == 0 {
		return nil, ErrNoGuests
	}

	out := make([]Rendered, 0, len(guests))

	for _, guest := range guests {
		invitation := NewInvitation(guest, baseURL, husband, wife)

		var buf bytes.Buffer
		if err := parsed.Execute(&buf, invitation); err != nil {
			return nil, fmt.Errorf("rendering the invite for %q failed: %w", guest.ID, err)
		}

		out = append(out, Rendered{Invitation: invitation, Message: buf.String()})
	}

	return out, nil
}

// QRCode renders a guest's link as a PNG.
func QRCode(link string) ([]byte, error) {
	png, err := qrcode.Encode(link, qrcode.Medium, qrSize)
	if err != nil {
		return nil, fmt.Errorf("cannot encode %q as a qr code: %w", link, err)
	}

	return png, nil
}
