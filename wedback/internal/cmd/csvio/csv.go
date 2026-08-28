// Package csvio implements bulk import and export of the guest list, so a
// wedding with more than a handful of guests does not have to go through the
// interactive `insert` form one person at a time.
package csvio

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
)

var (
	ErrMissingColumns = errors.New("missing required columns")
	ErrNoRows         = errors.New("file contains no guest rows")
)

// Columns lists the columns Parse understands, in the order Export writes them.
func Columns() []string {
	return []string{
		"first_name",
		"last_name",
		"spouse_first_name",
		"spouse_last_name",
		"is_family",
		"children",
	}
}

// exportColumns adds the read-only columns that only make sense on the way out.
func exportColumns() []string {
	out := append([]string{"id"}, Columns()...)

	return append(out, "coming", "plus_one", "answered")
}

// Row is one parsed CSV line, ready to hand to the guest service.
type Row struct {
	// Line is the 1-based line number in the source file, for error messages.
	Line            int
	FirstName       string
	LastName        string
	SpouseFirstName string
	SpouseLastName  string
	IsFamily        bool
	Children        int
}

// lookup resolves column names to their position in a record.
type lookup map[string]int

func (l lookup) get(record []string, name string) string {
	i, ok := l[name]
	if !ok || i >= len(record) {
		return ""
	}

	return strings.TrimSpace(record[i])
}

// Parse reads guests from r. Columns are matched by header name, so their
// order in the file does not matter and unknown columns are ignored.
func Parse(r io.Reader) ([]Row, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrNoRows
		}

		return nil, fmt.Errorf("cannot read csv header: %w", err)
	}

	columns, err := newLookup(header)
	if err != nil {
		return nil, err
	}

	var rows []Row

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("cannot read csv: %w", err)
		}

		// Records are read one at a time so FieldPos reports the true file
		// line: encoding/csv silently drops blank lines, so counting records
		// would drift away from what the author sees in their editor.
		line, _ := reader.FieldPos(0)

		// Skip comma-only lines rather than failing a whole import over them.
		if isBlank(record) {
			continue
		}

		row, err := parseRow(columns, record, line)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, ErrNoRows
	}

	return rows, nil
}

func newLookup(header []string) (lookup, error) {
	columns := make(lookup, len(header))
	for i, name := range header {
		columns[strings.ToLower(strings.TrimSpace(name))] = i
	}

	var missing []string

	for _, name := range []string{"first_name", "last_name"} {
		if _, ok := columns[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrMissingColumns, strings.Join(missing, ", "))
	}

	return columns, nil
}

func parseRow(columns lookup, record []string, line int) (Row, error) {
	children, err := parseChildren(columns.get(record, "children"), line)
	if err != nil {
		return Row{}, err // nolint: exhaustruct
	}

	isFamily, err := parseFamily(columns.get(record, "is_family"), line)
	if err != nil {
		return Row{}, err // nolint: exhaustruct
	}

	return Row{
		Line:            line,
		FirstName:       columns.get(record, "first_name"),
		LastName:        columns.get(record, "last_name"),
		SpouseFirstName: columns.get(record, "spouse_first_name"),
		SpouseLastName:  columns.get(record, "spouse_last_name"),
		IsFamily:        isFamily,
		Children:        children,
	}, nil
}

func parseChildren(raw string, line int) (int, error) {
	if raw == "" {
		return 0, nil
	}

	children, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("line %d: children must be a number: %w", line, err)
	}

	if children < 0 {
		return 0, fmt.Errorf("%w: line %d: children cannot be negative", ErrNegativeChildren, line)
	}

	return children, nil
}

// ErrNegativeChildren reports a guest row asking for a negative child count.
var ErrNegativeChildren = errors.New("invalid children count")

func parseFamily(raw string, line int) (bool, error) {
	if raw == "" {
		return false, nil
	}

	isFamily, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("line %d: is_family must be true or false: %w", line, err)
	}

	return isFamily, nil
}

// Export writes guests to w, including their RSVP answers. When baseURL is not
// empty an extra column carries each guest's personal invitation link.
func Export(w io.Writer, guests []model.Guest, baseURL string) error {
	writer := csv.NewWriter(w)

	header := exportColumns()
	if baseURL != "" {
		header = append(header, "link")
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("cannot write header: %w", err)
	}

	for _, guest := range guests {
		record := []string{
			guest.ID,
			guest.FirstName,
			guest.LastName,
			deref(guest.SpouseFirstName),
			deref(guest.SpouseLastName),
			strconv.FormatBool(guest.IsFamily),
			strconv.Itoa(guest.Children),
			strconv.FormatBool(guest.Coming()),
			strconv.FormatBool(guest.PlusOne()),
			strconv.FormatBool(guest.Answer != nil),
		}

		if baseURL != "" {
			record = append(record, strings.TrimSuffix(baseURL, "/")+"/guests/"+guest.ID)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("cannot write guest %q: %w", guest.ID, err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("cannot flush csv: %w", err)
	}

	return nil
}

func isBlank(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}

	return true
}

func deref(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
