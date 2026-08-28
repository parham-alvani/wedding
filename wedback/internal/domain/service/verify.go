package service

import (
	"context"
	"slices"
	"strings"
	"unicode"
)

// normaliseName folds a typed name into a comparable form: case, surrounding
// and repeated whitespace, and punctuation all stop mattering. Guests type
// their own name from memory, on a phone, so the comparison has to be
// forgiving without becoming meaningless.
func normaliseName(raw string) string {
	var b strings.Builder

	lastWasSpace := true

	for _, r := range strings.ToLower(strings.TrimSpace(raw)) {
		switch {
		case unicode.IsSpace(r):
			if !lastWasSpace {
				b.WriteRune(' ')

				lastWasSpace = true
			}
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)

			lastWasSpace = false
		default:
			// Punctuation and dashes are dropped rather than treated as a
			// separator, so "Alvani-Dastan" matches "alvani dastan" only via
			// the space rule above; here it becomes "alvanidastan".
		}
	}

	return strings.TrimSpace(b.String())
}

// acceptedNames lists every spelling that should let a guest in: either
// partner's first name, either last name, and the full names.
func acceptedNames(guest guestNames) []string {
	names := []string{guest.first, guest.last, guest.first + " " + guest.last}

	if guest.spouseFirst != "" {
		names = append(names, guest.spouseFirst)
	}

	if guest.spouseLast != "" {
		names = append(names, guest.spouseLast, guest.spouseFirst+" "+guest.spouseLast)
	}

	out := make([]string, 0, len(names))

	for _, name := range names {
		if normalised := normaliseName(name); normalised != "" {
			out = append(out, normalised)
		}
	}

	return out
}

type guestNames struct {
	first       string
	last        string
	spouseFirst string
	spouseLast  string
}

// VerifyName reports whether the given name belongs to the invitation. It is a
// soft gate: it keeps a forwarded link from opening for someone who does not
// know whose invitation it is, and it is not a substitute for a real secret.
func (svc GuestSvc) VerifyName(ctx context.Context, id string, name string) (bool, error) {
	guest, err := svc.Get(ctx, id)
	if err != nil {
		return false, err
	}

	typed := normaliseName(name)
	if typed == "" {
		return false, nil
	}

	names := guestNames{first: guest.FirstName, last: guest.LastName, spouseFirst: "", spouseLast: ""}
	if guest.SpouseFirstName != nil {
		names.spouseFirst = *guest.SpouseFirstName
	}

	if guest.SpouseLastName != nil {
		names.spouseLast = *guest.SpouseLastName
	}

	return slices.Contains(acceptedNames(names), typed), nil
}
