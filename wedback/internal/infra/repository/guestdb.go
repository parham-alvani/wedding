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
	"gorm.io/gorm/clause"
)

type GuestDB struct {
	g      gorm.Interface[model.Guest]
	a      gorm.Interface[model.Answer]
	logger *zap.Logger
}

func ProvideGuestDB(db *db.DB, logger *zap.Logger) *GuestDB {
	return &GuestDB{
		g:      gorm.G[model.Guest](db.DB),
		a:      gorm.G[model.Answer](db.DB),
		logger: logger.Named("repository.guestdb"),
	}
}

func (r *GuestDB) Create(ctx context.Context, guest model.Guest) error {
	if err := r.g.Create(ctx, &guest); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return guestrepo.ErrDuplicateGuestByName
		}

		r.logger.Error("guest creation failed", zap.Error(err), zap.String(logtag.Operation, "create"))

		return fmt.Errorf("guest creation failed %w", err)
	}

	return nil
}

func (r *GuestDB) Get(ctx context.Context, id string) (model.Guest, error) {
	guest, err := r.g.Joins(clause.LeftJoin.Association("Answer"), nil).Where("guests.id = ?", id).First(ctx)
	if err != nil {
		r.logger.Error("fetching guest from database failed", zap.Error(err), zap.String(logtag.Operation, "get"))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guest, guestrepo.ErrGuestNotFound
		}

		return guest, fmt.Errorf("fetching guest from database failed %w", err)
	}

	return guest, nil
}

func (r *GuestDB) List(ctx context.Context) ([]model.Guest, error) {
	guests, err := r.g.Joins(clause.LeftJoin.Association("Answer"), nil).Find(ctx)
	if err != nil {
		r.logger.Error("fetching guests from database failed", zap.Error(err), zap.String(logtag.Operation, "list"))

		return nil, fmt.Errorf("fetching guests from database failed %w", err)
	}

	return guests, nil
}

func (r *GuestDB) Update(ctx context.Context, guest model.Guest) error {
	if _, err := r.g.Where("id = ?", guest.ID).Updates(ctx, guest); err != nil {
		r.logger.Error("updating guest failed", zap.Error(err), zap.String(logtag.Operation, "update"))

		return fmt.Errorf("updating guest failed %w", err)
	}

	return nil
}

func (r *GuestDB) Answer(ctx context.Context, id string, answer model.Answer) error {
	guest, err := r.Get(ctx, id)
	if err != nil {
		r.logger.Error("guest fetching failed", zap.Error(err), zap.String(logtag.Operation, "answer"))

		return fmt.Errorf("guest fetching failed %w", err)
	}

	answer.GuestID = guest.ID

	// An existing reply is updated in place rather than inserted again: the
	// unique index on guest_id would reject a second row, and a guest who
	// changes their mind should not have to be reset first.
	if guest.Answer != nil {
		answer.ID = guest.Answer.ID

		// Select is required so that answering "false", or clearing a note,
		// is written: Updates skips zero-valued struct fields otherwise.
		if _, err := r.a.
			Where("guest_id = ?", guest.ID).
			Select("coming", "plus_one", "dietary", "song").
			Updates(ctx, answer); err != nil {
			r.logger.Error("answer update failed", zap.Error(err), zap.String(logtag.Operation, "answer"))

			return fmt.Errorf("answer update failed %w", err)
		}

		return nil
	}

	// nolint: gosec
	answer.ID = rand.Int64()

	if err := r.a.Create(ctx, &answer); err != nil {
		r.logger.Error("answer creation failed", zap.Error(err), zap.String(logtag.Operation, "answer"))

		return fmt.Errorf("answer creation failed %w", err)
	}

	return nil
}

func (r *GuestDB) ResetAnswer(ctx context.Context, id string) error {
	if _, err := r.g.Where("id = ?", id).First(ctx); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guestrepo.ErrGuestNotFound
		}

		return fmt.Errorf("guest fetching failed %w", err)
	}

	if _, err := r.a.Where("guest_id = ?", id).Delete(ctx); err != nil {
		r.logger.Error("answer reset failed", zap.Error(err), zap.String(logtag.Operation, "reset"))

		return fmt.Errorf("answer reset failed %w", err)
	}

	return nil
}
