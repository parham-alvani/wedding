package service

import (
	"context"
	"fmt"

	"github.com/parham-alvani/wedding/wedback/internal/domain/generator"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
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

func (svc GuestSvc) New(ctx context.Context, name string, partner string) (model.Guest, error) {
	guest := model.Guest{
		ID:     svc.generator.ID(),
		Name:   name,
		Spouse: partner,
		Answer: nil,
	}

	if err := svc.repository.Create(ctx, guest); err != nil {
		return guest, fmt.Errorf("guest creation failed %w", err)
	}

	return guest, nil
}

func (svc GuestSvc) Answer(ctx context.Context, id string, coming bool, plusOne bool) error {
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
