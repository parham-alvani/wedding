package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/parham-alvani/wedding/wedback/internal/domain/generator"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
)

var (
	ErrGuestNameRequired        = errors.New("first name and last name are required for a guest")
	ErrPartnerNameRequired      = errors.New("first name and last name are required for a guest's partner")
	ErrComingRequiredForPlusOne = errors.New("guest should come to have plus one")
	ErrRSVPClosed               = errors.New("the rsvp deadline has passed")
	ErrNoteTooLong              = errors.New("note is too long")
)

// MaxNoteLength caps the free-text fields a guest can submit. Invitation links
// are unauthenticated, so the fields are bounded on the way in.
const MaxNoteLength = 280

// Reply is everything a guest tells us when they answer.
type Reply struct {
	Coming  bool
	PlusOne bool
	// Dietary is what the kitchen needs to know; Song is a request for the
	// night. Both are optional free text.
	Dietary string
	Song    string
}

// normalise trims the free-text fields and rejects anything oversized.
func (r Reply) normalise() (Reply, error) {
	r.Dietary = strings.TrimSpace(r.Dietary)
	r.Song = strings.TrimSpace(r.Song)

	for field, value := range map[string]string{"dietary": r.Dietary, "song": r.Song} {
		if len(value) > MaxNoteLength {
			return r, fmt.Errorf("%w: %s is %d characters, the limit is %d",
				ErrNoteTooLong, field, len(value), MaxNoteLength)
		}
	}

	return r, nil
}

type GuestSvc struct {
	repository guestrepo.Repository
	generator  generator.Generator
	// deadline is the moment the RSVP closes. The zero time means it never
	// closes.
	deadline time.Time
}

func ProvideGuestSvc(
	repo guestrepo.Repository,
	gen generator.Generator,
	cfg wedding.Config,
) (GuestSvc, error) {
	svc := GuestSvc{
		repository: repo,
		generator:  gen,
		deadline:   time.Time{},
	}

	if raw := strings.TrimSpace(cfg.RSVPDeadline); raw != "" {
		deadline, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return svc, fmt.Errorf("wedding.rsvp_deadline %q is not an RFC 3339 timestamp: %w", raw, err)
		}

		svc.deadline = deadline
	}

	return svc, nil
}

// RSVPClosed reports whether the deadline has passed.
func (svc GuestSvc) RSVPClosed() bool {
	return !svc.deadline.IsZero() && time.Now().After(svc.deadline)
}

// Deadline returns the configured RSVP deadline, zero if there is none.
func (svc GuestSvc) Deadline() time.Time {
	return svc.deadline
}

// NewGuest is everything needed to add a guest to the list.
type NewGuest struct {
	FirstName        string
	LastName         string
	PartnerFirstName string
	PartnerLastName  string
	IsFamily         bool
	Children         int
	// Events is the ceremonies this guest is invited to, as a separated list
	// such as "engagement,wedding". Empty invites them to everything.
	Events string
}

func (svc GuestSvc) New(ctx context.Context, in NewGuest) (model.Guest, error) {
	fname := strings.TrimSpace(in.FirstName)
	lname := strings.TrimSpace(in.LastName)

	if len(lname) == 0 || len(fname) == 0 {
		return model.Guest{}, ErrGuestNameRequired // nolint: exhaustruct_v5
	}

	// Parsing normalises the order and rejects a typo before it is stored.
	events, err := model.ParseEvents(in.Events)
	if err != nil {
		return model.Guest{}, err // nolint: exhaustruct_v5, wrapcheck
	}

	guest := model.Guest{
		ID:              svc.generator.ID(),
		FirstName:       fname,
		LastName:        lname,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		IsFamily:        in.IsFamily,
		Children:        in.Children,
		Events:          model.FormatEvents(events),
	}

	partnerFname := strings.TrimSpace(in.PartnerFirstName)
	partnerLname := strings.TrimSpace(in.PartnerLastName)

	if (len(partnerFname) != 0) && (len(partnerLname) != 0) {
		guest.SpouseFirstName = &partnerFname
		guest.SpouseLastName = &partnerLname
	} else if (len(partnerFname) != 0) || (len(partnerLname) != 0) {
		return model.Guest{}, ErrPartnerNameRequired // nolint: exhaustruct_v5
	}

	if err := svc.repository.Create(ctx, guest); err != nil {
		return guest, fmt.Errorf("guest creation failed %w", err)
	}

	return guest, nil
}

func (svc GuestSvc) Answer(ctx context.Context, id string, reply Reply) error {
	if svc.RSVPClosed() {
		return ErrRSVPClosed
	}

	reply, err := reply.normalise()
	if err != nil {
		return err
	}

	if err := svc.repository.Answer(ctx, id, model.Answer{
		ID:      0,
		Coming:  reply.Coming,
		PlusOne: reply.PlusOne,
		Dietary: reply.Dietary,
		Song:    reply.Song,
		GuestID: "",
	}); err != nil {
		return fmt.Errorf("answer creation failed %w", err)
	}

	return nil
}

// ResetAnswer clears a guest's reply. This is an organiser action, so it is
// deliberately not gated on the RSVP deadline.
func (svc GuestSvc) ResetAnswer(ctx context.Context, id string) error {
	if err := svc.repository.ResetAnswer(ctx, id); err != nil {
		return fmt.Errorf("answer reset failed %w", err)
	}

	return nil
}

func (svc GuestSvc) Get(ctx context.Context, id string) (model.Guest, error) {
	guest, err := svc.repository.Get(ctx, id)
	if err != nil {
		return guest, fmt.Errorf("guest fetching failed %w", err)
	}

	return guest, nil
}

func (svc GuestSvc) List(ctx context.Context) ([]model.Guest, error) {
	guests, err := svc.repository.List(ctx)
	if err != nil {
		return guests, fmt.Errorf("guests fetching failed %w", err)
	}

	return guests, nil
}
