package csvio_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/csvio"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMinimalColumns(t *testing.T) {
	t.Parallel()

	rows, err := csvio.Parse(strings.NewReader("first_name,last_name\nAli,Irani\n"))
	require.NoError(t, err)
	require.Len(t, rows, 1)

	assert.Equal(t, "Ali", rows[0].FirstName)
	assert.Equal(t, "Irani", rows[0].LastName)
	assert.False(t, rows[0].IsFamily)
	assert.Equal(t, 0, rows[0].Children)
	assert.Equal(t, 2, rows[0].Line)
}

func TestParseIgnoresColumnOrderAndUnknownColumns(t *testing.T) {
	t.Parallel()

	rows, err := csvio.Parse(strings.NewReader(
		"children,note,last_name,first_name,is_family\n2,vip,Shirazi,Reza,true\n",
	))
	require.NoError(t, err)
	require.Len(t, rows, 1)

	assert.Equal(t, "Reza", rows[0].FirstName)
	assert.Equal(t, "Shirazi", rows[0].LastName)
	assert.True(t, rows[0].IsFamily)
	assert.Equal(t, 2, rows[0].Children)
}

func TestParseTrimsAndSkipsBlankLines(t *testing.T) {
	t.Parallel()

	rows, err := csvio.Parse(strings.NewReader(
		"first_name,last_name\n  Sara , Tehrani \n\n,\nNima,Karimi\n",
	))
	require.NoError(t, err)
	require.Len(t, rows, 2)

	assert.Equal(t, "Sara", rows[0].FirstName)
	assert.Equal(t, "Tehrani", rows[0].LastName)
	assert.Equal(t, "Nima", rows[1].FirstName)
}

// encoding/csv drops blank lines entirely, so line numbers have to come from
// the reader rather than from counting records.
func TestParseReportsTrueFileLineAfterBlankLines(t *testing.T) {
	t.Parallel()

	rows, err := csvio.Parse(strings.NewReader(
		"first_name,last_name\nAli,Irani\n\n\nNima,Karimi\n",
	))
	require.NoError(t, err)
	require.Len(t, rows, 2)

	assert.Equal(t, 2, rows[0].Line)
	assert.Equal(t, 5, rows[1].Line, "Nima is on line 5 of the file")
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		input string
		want  error
	}{
		"missing required columns": {input: "name,last_name\nAli,Irani\n", want: csvio.ErrMissingColumns},
		"header only":              {input: "first_name,last_name\n", want: csvio.ErrNoRows},
		"empty file":               {input: "", want: csvio.ErrNoRows},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := csvio.Parse(strings.NewReader(tc.input))
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func TestParseRejectsBadScalars(t *testing.T) {
	t.Parallel()

	for name, input := range map[string]string{
		"children not a number": "first_name,last_name,children\nAli,Irani,many\n",
		"children negative":     "first_name,last_name,children\nAli,Irani,-1\n",
		"is_family not a bool":  "first_name,last_name,is_family\nAli,Irani,maybe\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := csvio.Parse(strings.NewReader(input))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "line 2")
		})
	}
}

func TestExportRoundTrip(t *testing.T) {
	t.Parallel()

	spouseFirst, spouseLast := "Maryam", "Akhyani"
	guests := []model.Guest{
		{
			ID:              "aK3nQ7pLx2",
			FirstName:       "Ali",
			LastName:        "Irani",
			SpouseFirstName: &spouseFirst,
			SpouseLastName:  &spouseLast,
			IsFamily:        false,
			Children:        0,
			Answer:          &model.Answer{ID: 1, Coming: true, PlusOne: true, GuestID: "aK3nQ7pLx2"},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, csvio.Export(&buf, guests, "https://example.test"))

	out := buf.String()
	assert.Contains(t, out, "Ali,Irani,Maryam,Akhyani")
	assert.Contains(t, out, "https://example.test/guests/aK3nQ7pLx2")
	assert.Contains(t, out, "id,first_name,last_name")

	// The export must be re-importable.
	rows, err := csvio.Parse(strings.NewReader(out))
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "Ali", rows[0].FirstName)
	assert.Equal(t, "Maryam", rows[0].SpouseFirstName)
}

func TestExportWithoutLinks(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	guest := model.Guest{
		ID:              "x",
		FirstName:       "A",
		LastName:        "B",
		SpouseFirstName: nil,
		SpouseLastName:  nil,
		IsFamily:        false,
		Children:        0,
		Answer:          nil,
	}

	require.NoError(t, csvio.Export(&buf, []model.Guest{guest}, ""))

	assert.NotContains(t, buf.String(), "link")
}
