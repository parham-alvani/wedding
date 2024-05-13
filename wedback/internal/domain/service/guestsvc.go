package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/parham-alvani/wedding/wedback/internal/domain/generator"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
)

var (
	ErrGuestNameRequired        = errors.New("first name and last name are required for a guest")
	ErrPartnerNameRequired      = errors.New("first name and last name are required for a guest's partner")
	ErrComingRequiredForPlusOne = errors.New("guest should come to have plus one")
)

type GuestSvc struct {
	repository guestrepo.Repository
	generator  generator.Generator
}

func ProvideGuestSvc(repo guestrepo.Repository, gen generator.Generator) GuestSvc {
	return GuestSvc{
		repository: repo,
		generator:  gen,
	}
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
	if !coming && plusOne {
		return ErrComingRequiredForPlusOne
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
