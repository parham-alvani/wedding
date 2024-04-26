package repository

import (
	"context"
	"errors"
	"fmt"

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

	if err := r.db.DB.WithContext(ctx).Where("id = ?", id).First(&guest).Error; err != nil {
		r.logger.Error("fetching guest from database failed", zap.Error(err), zap.String(logtag.Operation, "get"))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guest, guestrepo.ErrGuestNotFound
		}

		return guest, fmt.Errorf("fetching guest from database failed %w", err)
	}

	return guest, nil
}

func (r *GuestDB) Update(ctx context.Context, guest model.Guest) error {
	if err := r.db.DB.WithContext(ctx).Save(guest).Error; err != nil {
		r.logger.Error("updating guest failed", zap.Error(err), zap.String(logtag.Operation, "update"))

		return fmt.Errorf("updating guest failed %w", err)
	}

	return nil
}
