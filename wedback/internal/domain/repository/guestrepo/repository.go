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
	Answer(ctx context.Context, id string, aswer model.Answer) error
}
