package guestrepo

import (
	"context"
	"errors"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
)

var (
	ErrGuestNotFound        = errors.New("guest does not exist")
	ErrDuplicateGuestByName = errors.New("guest with a same name already exists")
)

type Repository interface {
	Create(ctx context.Context, guest model.Guest) error
	Update(ctx context.Context, guest model.Guest) error
	Get(ctx context.Context, id string) (model.Guest, error)
	List(ctx context.Context) ([]model.Guest, error)
	// Answer records a guest's reply, replacing any earlier one so a guest
	// can change their mind.
	Answer(ctx context.Context, id string, answer model.Answer) error
	// ResetAnswer removes a guest's reply, putting them back in the waiting
	// list. Resetting a guest who never answered is not an error.
	ResetAnswer(ctx context.Context, id string) error
}
