package list

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/parham-alvani/wedding/wedback/internal/cmd/app"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

const (
	colorOrange     = lipgloss.Color("#FF8C00")
	colorSoftOrange = lipgloss.Color("#FFB347")
	colorDimText    = lipgloss.Color("#997040")
)

type guestsModel struct {
	repository guestrepo.Repository
	wedding    wedding.Config

	// onlyWaiting hides everyone who has already replied. Family guests count
	// as attending without replying, so they are never waiting.
	onlyWaiting bool

	isLoading bool
	width     int
	height    int

	table   table.Model
	spinner spinner.Model
	summary string
}

type guestsListMsg struct {
	guests []model.Guest
}

func (m guestsModel) List() tea.Msg {
	guests, err := m.repository.List(context.Background())
	if err != nil {
		pterm.Error.Printfln("reading from database failed %s", err)

		return nil
	}

	return guestsListMsg{
		guests: guests,
	}
}

func (m guestsModel) Init() tea.Cmd {
	return m.List
}

// nolint: funlen, cyclop
func (m guestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg == nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(max(m.height-10, 5)) //nolint: mnd
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case guestsListMsg:
		guests := m.filter(msg.guests)
		rows := make([]table.Row, len(guests))

		comingGuests := 0
		notComingGuests := 0
		notAnsweredGuests := 0

		for i, guest := range guests {
			spouseFirstName := ""
			if guest.SpouseFirstName != nil {
				spouseFirstName = *guest.SpouseFirstName
			}

			spouseLastName := ""
			if guest.SpouseLastName != nil {
				spouseLastName = *guest.SpouseLastName
			}

			// nolint: nestif, mnd
			if !guest.IsFamily && guest.Answer == nil {
				notAnsweredGuests++
			} else {
				if guest.Coming() {
					if guest.PlusOne() {
						comingGuests += (2 + guest.Children)
					} else {
						comingGuests += (1 + guest.Children)
					}
				} else {
					notComingGuests++
				}
			}

			rows[i] = table.Row{
				guest.FirstName,
				guest.LastName,
				spouseFirstName,
				spouseLastName,
				strconv.FormatBool(guest.IsFamily),
				strconv.Itoa(guest.Children),
				guest.ID,
				formatBool(guest.Coming()),
				formatBool(guest.PlusOne()),
				formatBool(guest.Answer == nil && !guest.IsFamily),
			}
		}

		m.isLoading = false
		m.table.SetRows(rows)

		m.summary = fmt.Sprintf(
			"Total: %d  |  Coming: %d  |  Not Coming: %d  |  Waiting: %d",
			len(guests),
			comingGuests,
			notComingGuests,
			notAnsweredGuests,
		)
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m guestsModel) View() string {
	if m.isLoading {
		return fmt.Sprintf("\n\n   %s Loading guests...\n\n", m.spinner.View())
	}

	summaryStyle := lipgloss.NewStyle().
		Foreground(colorSoftOrange).
		Bold(true).
		Padding(1, 0) //nolint: mnd

	helpStyle := lipgloss.NewStyle().
		Foreground(colorDimText).
		Italic(true)

	return m.table.View() + "\n" +
		summaryStyle.Render(m.summary) + "\n" +
		helpStyle.Render("↑/↓ navigate • q quit")
}

// filter applies the --waiting flag. Family guests count as attending without
// replying, so they are never waiting.
func (m guestsModel) filter(guests []model.Guest) []model.Guest {
	if !m.onlyWaiting {
		return guests
	}

	waiting := make([]model.Guest, 0, len(guests))

	for _, guest := range guests {
		if guest.Answer == nil && !guest.IsFamily {
			waiting = append(waiting, guest)
		}
	}

	return waiting
}

func formatBool(b bool) string {
	if b {
		return "✓"
	}

	return "✗"
}

func weddingTableStyles() table.Styles {
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderForeground(colorOrange).
		Foreground(colorOrange).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#FF8C00")).
		Bold(false)

	s.Cell = s.Cell.
		Foreground(lipgloss.Color("#DDD"))

	return s
}

// nolint: mnd
func main(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	repository guestrepo.Repository,
	weddingCfg wedding.Config,
	onlyWaiting bool,
) {
	columns := []table.Column{
		{Title: "First Name", Width: 14},
		{Title: "Last Name", Width: 14},
		{Title: "Spouse First", Width: 14},
		{Title: "Spouse Last", Width: 14},
		{Title: "Family", Width: 7},
		{Title: "Kids", Width: 5},
		{Title: "ID", Width: 12},
		{Title: "Coming", Width: 7},
		{Title: "+1", Width: 4},
		{Title: "Waiting", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	t.SetStyles(weddingTableStyles())

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(colorOrange)

	dm := guestsModel{
		repository:  repository,
		wedding:     weddingCfg,
		onlyWaiting: onlyWaiting,
		isLoading:   true,
		width:       0,
		height:      0,
		spinner:     s,
		table:       t,
		summary:     "",
	}

	p := tea.NewProgram(dm, tea.WithAltScreen())

	lc.Append(
		fx.StartHook(func(_ context.Context) error {
			go func() {
				if _, err := p.Run(); err != nil {
					pterm.Fatal.Printfln("execution failed %s", err)
				}

				_ = shutdowner.Shutdown()
			}()

			return nil
		}),
	)
}

// Register list command.
func Register() *cli.Command {
	//nolint: exhaustruct_v5
	return &cli.Command{
		Name:        "list",
		Description: "List all guests with their RSVP status",
		Flags: []cli.Flag{
			//nolint: exhaustruct_v5
			&cli.BoolFlag{
				Name:  "waiting",
				Usage: "only show guests who have not replied yet",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			onlyWaiting := cmd.Bool("waiting")

			return app.Run(
				app.Providers(),
				fx.Supply(onlyWaiting),
				fx.Invoke(main),
			)
		},
	}
}
