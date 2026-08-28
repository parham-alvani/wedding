package service_test

import (
	"context"
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFirst  = "Ali"
	testLast   = "Irani"
	testConfig = "simple"
)

// oneGuestRepo returns the same guest for any id.
type oneGuestRepo struct {
	stubRepo

	guest model.Guest
}

func (r *oneGuestRepo) Get(context.Context, string) (model.Guest, error) {
	return r.guest, nil
}

func verifySvc(t *testing.T, guest model.Guest) service.GuestSvc {
	t.Helper()

	repo := &oneGuestRepo{stubRepo{last: model.Answer{}, calls: 0}, guest} // nolint: exhaustruct_v5
	svc, err := service.ProvideGuestSvc(repo, generator.Provide(generator.Config{Type: testConfig}), openConfig())
	require.NoError(t, err)

	return svc
}

func TestVerifyNameAcceptsEitherPartner(t *testing.T) {
	t.Parallel()

	spouseFirst, spouseLast := "Maryam", "Akhyani"
	svc := verifySvc(t, model.Guest{ // nolint: exhaustruct_v5
		ID: "x", FirstName: testFirst, LastName: testLast,
		SpouseFirstName: &spouseFirst, SpouseLastName: &spouseLast,
	})

	for _, name := range []string{
		testFirst, testLast, testFirst + " " + testLast, "Maryam", "Akhyani", "Maryam Akhyani",
		// Guests type from memory, on a phone.
		"  ali  ", "IRANI", "ali   irani", "Ali\tIrani",
	} {
		ok, err := svc.VerifyName(context.Background(), "x", name)
		require.NoError(t, err)
		assert.True(t, ok, "%q should open the invitation", name)
	}
}

func TestVerifyNameRejectsOthers(t *testing.T) {
	t.Parallel()

	svc := verifySvc(t, model.Guest{ // nolint: exhaustruct_v5
		ID: "x", FirstName: testFirst, LastName: testLast,
	})

	for _, name := range []string{"", "   ", "Sara", "Tehrani", "Ali Tehrani", "Al", "Iranian"} {
		ok, err := svc.VerifyName(context.Background(), "x", name)
		require.NoError(t, err)
		assert.False(t, ok, "%q should not open the invitation", name)
	}
}

func TestVerifyNameIgnoresPunctuation(t *testing.T) {
	t.Parallel()

	svc := verifySvc(t, model.Guest{ // nolint: exhaustruct_v5
		ID: "x", FirstName: testFirst, LastName: "Alvani-Dastan",
	})

	ok, err := svc.VerifyName(context.Background(), "x", "alvanidastan")
	require.NoError(t, err)
	assert.True(t, ok)
}

// A guest with no partner must not be opened by an empty spouse field.
func TestVerifyNameWithoutSpouse(t *testing.T) {
	t.Parallel()

	svc := verifySvc(t, model.Guest{ // nolint: exhaustruct_v5
		ID: "x", FirstName: "Kobra", LastName: "Dastan",
	})

	ok, err := svc.VerifyName(context.Background(), "x", " ")
	require.NoError(t, err)
	assert.False(t, ok)
}
