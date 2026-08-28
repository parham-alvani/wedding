package invite_test

import (
	"strings"
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/invite"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func guest(id, first, last string, spouse *string, family bool, answer *model.Answer) model.Guest {
	return model.Guest{
		ID:              id,
		FirstName:       first,
		LastName:        last,
		SpouseFirstName: spouse,
		SpouseLastName:  nil,
		IsFamily:        family,
		Children:        0,
		Answer:          answer,
	}
}

func TestLinkJoinsWithoutDoubleSlash(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "https://a.test/guests/xyz", invite.Link("https://a.test", "xyz"))
	assert.Equal(t, "https://a.test/guests/xyz", invite.Link("https://a.test/", "xyz"))
}

func TestNamesCombineCouples(t *testing.T) {
	t.Parallel()

	maryam := "Maryam"

	solo := invite.NewInvitation(guest("a", "Ali", "Irani", nil, false, nil), "https://a.test", "P", "E")
	assert.Equal(t, "Ali", solo.Names)

	couple := invite.NewInvitation(guest("b", "Ali", "Irani", &maryam, false, nil), "https://a.test", "P", "E")
	assert.Equal(t, "Ali & Maryam", couple.Names)
}

// Family guests are treated as attending without replying, so chasing them
// would be wrong.
func TestWaitingExcludesFamilyAndAnswered(t *testing.T) {
	t.Parallel()

	answered := &model.Answer{ID: 1, Coming: true, PlusOne: false, GuestID: "c"}

	guests := []model.Guest{
		guest("a", "Sara", "Tehrani", nil, false, nil),   // waiting
		guest("b", "Reza", "Shirazi", nil, true, nil),    // family
		guest("c", "Ali", "Irani", nil, false, answered), // replied
		guest("d", "Nima", "Karimi", nil, false, nil),    // waiting
	}

	waiting := invite.Waiting(guests)
	require.Len(t, waiting, 2)
	assert.Equal(t, "Sara", waiting[0].FirstName)
	assert.Equal(t, "Nima", waiting[1].FirstName)
}

func TestRenderUsesDefaultTemplateWhenBlank(t *testing.T) {
	t.Parallel()

	out, err := invite.Render("  ", []model.Guest{guest("a", "Ali", "Irani", nil, false, nil)},
		"https://a.test", "Parham", "Elaheh")
	require.NoError(t, err)
	require.Len(t, out, 1)

	assert.Contains(t, out[0].Message, "Dear Ali,")
	assert.Contains(t, out[0].Message, "Elaheh and Parham")
	assert.Contains(t, out[0].Message, "https://a.test/guests/a")
}

func TestRenderCustomTemplate(t *testing.T) {
	t.Parallel()

	out, err := invite.Render("{{ .Names }} -> {{ .Link }}",
		[]model.Guest{guest("zz", "Sara", "Tehrani", nil, false, nil)},
		"https://a.test/", "Parham", "Elaheh")
	require.NoError(t, err)
	assert.Equal(t, "Sara -> https://a.test/guests/zz", out[0].Message)
}

func TestRenderRejectsBadTemplate(t *testing.T) {
	t.Parallel()

	_, err := invite.Render("{{ .Names ", []model.Guest{guest("a", "Ali", "Irani", nil, false, nil)},
		"https://a.test", "P", "E")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invite template is invalid")
}

func TestRenderRejectsUnknownField(t *testing.T) {
	t.Parallel()

	_, err := invite.Render("{{ .Nope }}", []model.Guest{guest("a", "Ali", "Irani", nil, false, nil)},
		"https://a.test", "P", "E")
	require.Error(t, err)
}

func TestRenderNoGuests(t *testing.T) {
	t.Parallel()

	_, err := invite.Render("", nil, "https://a.test", "P", "E")
	require.ErrorIs(t, err, invite.ErrNoGuests)
}

func TestQRCodeIsAPNG(t *testing.T) {
	t.Parallel()

	png, err := invite.QRCode("https://a.test/guests/abc")
	require.NoError(t, err)
	require.NotEmpty(t, png)

	assert.True(t, strings.HasPrefix(string(png[:8]), "\x89PNG"), "should carry the PNG magic bytes")
}
