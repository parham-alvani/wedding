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

func (svc GuestSvc) New(ctx context.Context, name string) (model.Guest, error) {
	guest := model.Guest{
		ID:     svc.generator.ID(),
		Name:   name,
		Answer: nil,
	}

	if err := svc.repository.Create(ctx, guest); err != nil {
		return guest, fmt.Errorf("guest creation failed %w", err)
	}

	return guest, nil
}
