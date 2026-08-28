package model_test

import (
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An empty list is how an untiered guest list is stored, so it has to mean
// "invited to everything" rather than "invited to nothing".
func TestParseEventsEmptyMeansAll(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{"", "   ", "\t"} {
		events, err := model.ParseEvents(raw)
		require.NoError(t, err)
		assert.Equal(t, model.AllEvents(), events)
	}
}

func TestParseEventsAcceptsSeveralSeparators(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		"engagement,wedding", "engagement wedding", "engagement|wedding",
		"wedding, engagement", " WEDDING , Engagement ",
	} {
		events, err := model.ParseEvents(raw)
		require.NoError(t, err, raw)
		// Order is normalised to the order the ceremonies happen in.
		assert.Equal(t, []model.Event{model.EventEngagement, model.EventWedding}, events, raw)
	}
}

func TestParseEventsSingle(t *testing.T) {
	t.Parallel()

	events, err := model.ParseEvents("wedding")
	require.NoError(t, err)
	assert.Equal(t, []model.Event{model.EventWedding}, events)
}

func TestParseEventsDeduplicates(t *testing.T) {
	t.Parallel()

	events, err := model.ParseEvents("wedding,wedding,wedding")
	require.NoError(t, err)
	assert.Equal(t, []model.Event{model.EventWedding}, events)
}

func TestParseEventsRejectsUnknown(t *testing.T) {
	t.Parallel()

	_, err := model.ParseEvents("wedding,brunch")
	require.ErrorIs(t, err, model.ErrUnknownEvent)
	assert.Contains(t, err.Error(), "brunch")
}

func TestFormatEventsIsStable(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "engagement,wedding",
		model.FormatEvents([]model.Event{model.EventWedding, model.EventEngagement}))
	assert.Equal(t, "wedding", model.FormatEvents([]model.Event{model.EventWedding}))
}

func TestGuestInvitedTo(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		events     string
		engagement bool
		wedding    bool
	}{
		"untiered means both": {events: "", engagement: true, wedding: true},
		"wedding only":        {events: "wedding", engagement: false, wedding: true},
		"engagement only":     {events: "engagement", engagement: true, wedding: false},
		"both":                {events: "engagement,wedding", engagement: true, wedding: true},
		// A value written outside the service should not silently uninvite
		// somebody; fall back to inviting them to everything.
		"corrupt falls back": {events: "brunch", engagement: true, wedding: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			guest := model.Guest{ // nolint: exhaustruct_v5
				ID: "x", FirstName: "A", LastName: "B", Events: tc.events,
			}

			assert.Equal(t, tc.engagement, guest.InvitedTo(model.EventEngagement))
			assert.Equal(t, tc.wedding, guest.InvitedTo(model.EventWedding))
		})
	}
}
