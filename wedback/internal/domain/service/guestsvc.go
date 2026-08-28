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
)

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

func (svc GuestSvc) New(
	ctx context.Context,
	fname string,
	lname string,
	partnerFname string,
	partnerLname string,
	isFamily bool,
	children int,
) (model.Guest, error) {
	fname = strings.TrimSpace(fname)
	lname = strings.TrimSpace(lname)

	if len(lname) == 0 || len(fname) == 0 {
		return model.Guest{}, ErrGuestNameRequired
	}

	guest := model.Guest{
		ID:              svc.generator.ID(),
		FirstName:       fname,
		LastName:        lname,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		IsFamily:        isFamily,
		Children:        children,
	}

	partnerFname = strings.TrimSpace(partnerFname)
	partnerLname = strings.TrimSpace(partnerLname)

	if (len(partnerFname) != 0) && (len(partnerLname) != 0) {
		guest.SpouseFirstName = &partnerFname
		guest.SpouseLastName = &partnerLname
	} else if (len(partnerFname) != 0) || (len(partnerLname) != 0) {
		return model.Guest{}, ErrPartnerNameRequired
	}

	if err := svc.repository.Create(ctx, guest); err != nil {
		return guest, fmt.Errorf("guest creation failed %w", err)
	}

	return guest, nil
}

func (svc GuestSvc) Answer(ctx context.Context, id string, coming bool, plusOne bool) error {
	if svc.RSVPClosed() {
		return ErrRSVPClosed
	}

	if err := svc.repository.Answer(ctx, id, model.Answer{
		ID:      0,
		Coming:  coming,
		PlusOne: plusOne,
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
