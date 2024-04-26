package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

type GuestDBTestSuite struct {
	suite.Suite

	repo guestrepo.Repository
	db   *db.DB

	app *fxtest.App
}

func (s *GuestDBTestSuite) SetupSuite() {
	s.app = fxtest.New(s.T(),
		fx.Provide(config.Provide),
		fx.Provide(logger.Provide),
		fx.Provide(db.Provide),
		fx.Provide(
			fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
		),
		fx.Invoke(func(repo guestrepo.Repository, db *db.DB) {
			s.db = db
			s.repo = repo
		}),
	).RequireStart()
}

func (s *GuestDBTestSuite) TearDownTest() {
	require := s.Require()

	// nolint: exhaustruct
	stmt := &gorm.Statement{DB: s.db.DB}
	require.NoError(stmt.Parse(new(model.Guest)))

	tx := s.db.DB.Exec(fmt.Sprintf("DELETE FROM %s;", stmt.Schema.Table))
	require.NoError(tx.Error)
}

func (s *GuestDBTestSuite) TearDownSuite() {
	s.app.RequireStop()
}

func (s *GuestDBTestSuite) TestNotFound() {
	require := s.Require()

	_, err := s.repo.Get(context.Background(), "static_random")
	require.ErrorIs(guestrepo.ErrGuestNotFound, err)
}

func (s *GuestDBTestSuite) TestCreate() {
	require := s.Require()

	// nolint: exhaustruct
	require.NoError(s.repo.Create(context.Background(), model.Guest{
		ID:     "unique",
		Name:   "Ali Irani",
		Answer: nil,
	}))

	guest, err := s.repo.Get(context.Background(), "unique")
	require.NoError(err)

	require.Equal("Ali Irani", guest.Name)
}

func TestGuestDB(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GuestDBTestSuite))
}
