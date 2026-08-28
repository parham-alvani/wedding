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

// Fixture values shared across the suite's cases.
const (
	testGuestID   = "unique"
	testFirstName = "Ali"
	testLastName  = "Irani"
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

	// nolint: exhaustruct_v5
	stmt := &gorm.Statement{DB: s.db.DB}
	{
		require.NoError(stmt.Parse(new(model.Guest)))

		tx := s.db.DB.Exec(fmt.Sprintf("DELETE FROM %s;", stmt.Schema.Table))
		require.NoError(tx.Error)
	}

	{
		require.NoError(stmt.Parse(new(model.Answer)))

		tx := s.db.DB.Exec(fmt.Sprintf("DELETE FROM %s;", stmt.Schema.Table))
		require.NoError(tx.Error)
	}
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

	// nolint: exhaustruct_v5
	require.NoError(s.repo.Create(context.Background(), model.Guest{
		ID:              testGuestID,
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		Events:          "",
	}))

	guest, err := s.repo.Get(context.Background(), testGuestID)
	require.NoError(err)

	require.Equal("Ali Irani", guest.FirstName+" "+guest.LastName)
}

func (s *GuestDBTestSuite) TestCreateWithAnswer() {
	require := s.Require()

	// nolint: exhaustruct_v5
	require.NoError(s.repo.Create(context.Background(), model.Guest{
		ID:              testGuestID,
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		Events:          "",
	}))

	require.NoError(s.repo.Answer(context.Background(), testGuestID, model.Answer{
		ID:      0,
		GuestID: "",
		PlusOne: true,
		Coming:  true,
		Dietary: "",
		Song:    "",
	}))

	guest, err := s.repo.Get(context.Background(), testGuestID)
	require.NoError(err)

	require.Equal("Ali Irani", guest.FirstName+" "+guest.LastName)
	require.NotNil(guest.Answer)
	require.True(guest.Coming())
	require.True(guest.PlusOne())
}

// A guest who changes their mind must overwrite their reply rather than
// insert a second one, which the unique index on guest_id would reject.
func (s *GuestDBTestSuite) TestAnswerIsReplacedNotDuplicated() {
	require := s.Require()
	ctx := context.Background()

	// nolint: exhaustruct_v5
	require.NoError(s.repo.Create(ctx, model.Guest{
		ID:              testGuestID,
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		Events:          "",
	}))

	require.NoError(s.repo.Answer(ctx, testGuestID, model.Answer{
		ID: 0, GuestID: "", Coming: true, PlusOne: true, Dietary: "", Song: "",
	}))

	// Answering "not coming" writes zero values, which a naive struct update
	// would silently skip.
	require.NoError(s.repo.Answer(ctx, testGuestID, model.Answer{
		ID: 0, GuestID: "", Coming: false, PlusOne: false, Dietary: "", Song: "",
	}))

	guest, err := s.repo.Get(ctx, testGuestID)
	require.NoError(err)
	require.NotNil(guest.Answer)
	require.False(guest.Answer.Coming)
	require.False(guest.Answer.PlusOne)

	var answers int64
	require.NoError(s.db.DB.Table("answers").Where("guest_id = ?", testGuestID).Count(&answers).Error)
	require.EqualValues(1, answers, "a changed answer must reuse the same row")

	// And back again.
	require.NoError(s.repo.Answer(ctx, testGuestID, model.Answer{
		ID: 0, GuestID: "", Coming: true, PlusOne: false, Dietary: "", Song: "",
	}))

	guest, err = s.repo.Get(ctx, testGuestID)
	require.NoError(err)
	require.True(guest.Answer.Coming)
	require.False(guest.Answer.PlusOne)
}

func (s *GuestDBTestSuite) TestResetAnswer() {
	require := s.Require()
	ctx := context.Background()

	// nolint: exhaustruct_v5
	require.NoError(s.repo.Create(ctx, model.Guest{
		ID:              testGuestID,
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		Answer:          nil,
		Events:          "",
	}))

	require.NoError(s.repo.Answer(ctx, testGuestID, model.Answer{
		ID: 0, GuestID: "", Coming: true, PlusOne: true, Dietary: "", Song: "",
	}))

	require.NoError(s.repo.ResetAnswer(ctx, testGuestID))

	guest, err := s.repo.Get(ctx, testGuestID)
	require.NoError(err)
	require.Nil(guest.Answer, "the guest is back in the waiting list")

	// Resetting twice is not an error.
	require.NoError(s.repo.ResetAnswer(ctx, testGuestID))
}

func (s *GuestDBTestSuite) TestResetAnswerUnknownGuest() {
	require := s.Require()

	require.ErrorIs(s.repo.ResetAnswer(context.Background(), "nope"), guestrepo.ErrGuestNotFound)
}

func (s *GuestDBTestSuite) TestUpdate() {
	require := s.Require()

	require.NoError(s.repo.Create(context.Background(), model.Guest{
		ID:              "update-me",
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		IsFamily:        false,
		Children:        0,
		Answer:          nil,
		Events:          "",
	}))

	spouseFirst := "Maryam"
	spouseLast := "Akhyani"

	require.NoError(s.repo.Update(context.Background(), model.Guest{
		ID:              "update-me",
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: &spouseFirst,
		SpouseLastName:  &spouseLast,
		IsFamily:        true,
		Children:        2,
		Answer:          nil,
		Events:          "",
	}))

	guest, err := s.repo.Get(context.Background(), "update-me")
	require.NoError(err)

	require.True(guest.IsFamily)
	require.Equal(2, guest.Children)
	require.NotNil(guest.SpouseFirstName)
	require.Equal("Maryam", *guest.SpouseFirstName)
}

func (s *GuestDBTestSuite) TestCreateWithDuplicateName() {
	require := s.Require()

	require.NoError(s.repo.Create(context.Background(), model.Guest{
		ID:              testGuestID,
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		IsFamily:        false,
		Children:        0,
		Answer:          nil,
		Events:          "",
	}))

	require.ErrorIs(s.repo.Create(context.Background(), model.Guest{
		ID:              "not-unique",
		FirstName:       testFirstName,
		LastName:        testLastName,
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		IsFamily:        false,
		Children:        0,
		Answer:          nil,
		Events:          "",
	}), guestrepo.ErrDuplicateGuestByName)
}

func (s *GuestDBTestSuite) TestCreateWithAnswerButWithoutGuest() {
	require := s.Require()

	require.Error(s.repo.Answer(context.Background(), "not-found", model.Answer{
		ID:      0,
		GuestID: "",
		PlusOne: true,
		Coming:  true,
		Dietary: "",
		Song:    "",
	}))
}

func (s *GuestDBTestSuite) TestList() {
	require := s.Require()

	// nolint: exhaustruct_v5
	for i := range 10 {
		require.NoError(s.repo.Create(context.Background(), model.Guest{
			ID:        fmt.Sprintf("unique %d", i),
			FirstName: testFirstName,
			LastName:  fmt.Sprintf("Irani %d", i),
			Answer:    nil,
			Events:    "",
		}))
	}

	guests, err := s.repo.List(context.Background())
	require.NoError(err)

	require.Len(guests, 10)
}

func TestGuestDB(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GuestDBTestSuite))
}
