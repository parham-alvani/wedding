package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound     = errors.New("entity does not exist")
	ErrDuplicateKey = errors.New("entity with the same key already exists")
)

// Repository is a generic interface for basic CRUD operations.
type Repository[T any] interface {
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Get(ctx context.Context, id string) (T, error)
	List(ctx context.Context) ([]T, error)
}
