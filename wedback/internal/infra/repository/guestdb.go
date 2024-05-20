package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logtag"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GuestDB struct {
	db     *db.DB
	logger *zap.Logger
}

func ProvideGuestDB(db *db.DB, logger *zap.Logger) *GuestDB {
	return &GuestDB{
		db:     db,
		logger: logger.Named("repository.guestdb"),
	}
}

func (r *GuestDB) Create(ctx context.Context, guest model.Guest) error {
	if err := r.db.DB.WithContext(ctx).Save(&guest).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return guestrepo.ErrDuplicateGuestByName
		}

		r.logger.Error("guest creation failed", zap.Error(err), zap.String(logtag.Operation, "create"))

		return fmt.Errorf("guest creation failed %w", err)
	}

	return nil
}

func (r *GuestDB) Get(ctx context.Context, id string) (model.Guest, error) {
	var guest model.Guest

	if err := r.db.DB.WithContext(ctx).Where("guests.id = ?", id).Joins("Answer").First(&guest).Error; err != nil {
		r.logger.Error("fetching guest from database failed", zap.Error(err), zap.String(logtag.Operation, "get"))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guest, guestrepo.ErrGuestNotFound
		}

		return guest, fmt.Errorf("fetching guest from database failed %w", err)
	}

	return guest, nil
}

func (r *GuestDB) List(ctx context.Context) ([]model.Guest, error) {
	var guests []model.Guest

	if err := r.db.DB.WithContext(ctx).Joins("Answer").Find(&guests).Error; err != nil {
		r.logger.Error("fetching guests from database failed", zap.Error(err), zap.String(logtag.Operation, "list"))

		return nil, fmt.Errorf("fetching guests from database failed %w", err)
	}

	return guests, nil
}

func (r *GuestDB) Update(ctx context.Context, guest model.Guest) error {
	if err := r.db.DB.WithContext(ctx).Save(guest).Error; err != nil {
		r.logger.Error("updating guest failed", zap.Error(err), zap.String(logtag.Operation, "update"))

		return fmt.Errorf("updating guest failed %w", err)
	}

	return nil
}

func (r *GuestDB) Answer(ctx context.Context, id string, answer model.Answer) error {
	// nolint: gosec
	answer.ID = rand.Int64()

	guest, err := r.Get(ctx, id)
	if err != nil {
		r.logger.Error("guest fetching failed", zap.Error(err), zap.String(logtag.Operation, "answer"))

		return fmt.Errorf("guest fetching failed %w", err)
	}

	guest.Answer = &answer

	if err := r.db.DB.WithContext(ctx).Updates(&guest).Error; err != nil {
		r.logger.Error("answer creation failed", zap.Error(err), zap.String(logtag.Operation, "answer"))

		return fmt.Errorf("answer creation failed %w", err)
	}

	return nil
}
