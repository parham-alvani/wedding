package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubRepo records the last answer it was handed.
type stubRepo struct {
	last  model.Answer
	calls int
}

func (s *stubRepo) Create(context.Context, model.Guest) error { return nil }
func (s *stubRepo) Update(context.Context, model.Guest) error { return nil }
func (s *stubRepo) Get(context.Context, string) (model.Guest, error) {
	return model.Guest{}, nil // nolint: exhaustruct_v5
}
func (s *stubRepo) List(context.Context) ([]model.Guest, error) { return nil, nil }
func (s *stubRepo) ResetAnswer(context.Context, string) error   { return nil }

func (s *stubRepo) Answer(_ context.Context, _ string, answer model.Answer) error {
	s.last = answer
	s.calls++

	return nil
}

func newSvc(t *testing.T, cfg wedding.Config) (service.GuestSvc, *stubRepo) {
	t.Helper()

	repo := &stubRepo{last: model.Answer{}, calls: 0} // nolint: exhaustruct_v5
	svc, err := service.ProvideGuestSvc(repo, generator.Provide(generator.Config{Type: "simple"}), cfg)
	require.NoError(t, err)

	return svc, repo
}

func openConfig() wedding.Config {
	return wedding.Config{
		HusbandName: "P", WifeName: "E", BaseURL: "https://a.test",
		RSVPDeadline: "", InviteTemplate: "",
	}
}

func TestAnswerTrimsNotes(t *testing.T) {
	t.Parallel()

	svc, repo := newSvc(t, openConfig())

	require.NoError(t, svc.Answer(context.Background(), "id", service.Reply{
		Coming: true, PlusOne: false,
		Dietary: "  vegetarian  ", Song: "\tEbi — Shab Nashini\n",
	}))

	assert.Equal(t, "vegetarian", repo.last.Dietary)
	assert.Equal(t, "Ebi — Shab Nashini", repo.last.Song)
}

func TestAnswerRejectsOversizedNotes(t *testing.T) {
	t.Parallel()

	for name, reply := range map[string]service.Reply{
		"dietary": {Coming: true, PlusOne: false, Dietary: strings.Repeat("x", service.MaxNoteLength+1), Song: ""},
		"song":    {Coming: true, PlusOne: false, Dietary: "", Song: strings.Repeat("x", service.MaxNoteLength+1)},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svc, repo := newSvc(t, openConfig())

			err := svc.Answer(context.Background(), "id", reply)
			require.ErrorIs(t, err, service.ErrNoteTooLong)
			assert.Zero(t, repo.calls, "nothing should reach the repository")
		})
	}
}

// The limit is counted in bytes, so a note at exactly the cap is accepted.
func TestAnswerAcceptsNoteAtTheLimit(t *testing.T) {
	t.Parallel()

	svc, repo := newSvc(t, openConfig())

	require.NoError(t, svc.Answer(context.Background(), "id", service.Reply{
		Coming: true, PlusOne: false,
		Dietary: strings.Repeat("x", service.MaxNoteLength), Song: "",
	}))
	assert.Equal(t, 1, repo.calls)
}

func TestAnswerRefusedAfterDeadline(t *testing.T) {
	t.Parallel()

	cfg := openConfig()
	cfg.RSVPDeadline = "2020-01-01T00:00:00Z"

	svc, repo := newSvc(t, cfg)

	assert.True(t, svc.RSVPClosed())
	require.ErrorIs(t,
		svc.Answer(context.Background(), "id", service.Reply{Coming: true, PlusOne: false, Dietary: "", Song: ""}),
		service.ErrRSVPClosed)
	assert.Zero(t, repo.calls)
}

func TestAnswerAllowedBeforeDeadline(t *testing.T) {
	t.Parallel()

	cfg := openConfig()
	cfg.RSVPDeadline = "2100-01-01T00:00:00Z"

	svc, repo := newSvc(t, cfg)

	assert.False(t, svc.RSVPClosed())
	require.NoError(t,
		svc.Answer(context.Background(), "id", service.Reply{Coming: true, PlusOne: false, Dietary: "", Song: ""}))
	assert.Equal(t, 1, repo.calls)
}

func TestMalformedDeadlineIsRejected(t *testing.T) {
	t.Parallel()

	cfg := openConfig()
	cfg.RSVPDeadline = "next tuesday"

	repo := &stubRepo{last: model.Answer{}, calls: 0} // nolint: exhaustruct_v5
	_, err := service.ProvideGuestSvc(repo, generator.Provide(generator.Config{Type: "simple"}), cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RFC 3339")
}
