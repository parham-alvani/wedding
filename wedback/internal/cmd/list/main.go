package list

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

type guestsModel struct {
	repository guestrepo.Repository

	isLoading bool

	table   table.Model
	spinner spinner.Model
	text    textarea.Model
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
	case tea.KeyMsg:
		// nolint: exhaustive
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyDelete, tea.KeyBackspace:
			return m, tea.Batch(
				tea.Printf("Let's not having %s!", m.table.SelectedRow()[0]),
			)
		}
	case guestsListMsg:
		rows := make([]table.Row, len(msg.guests))

		comingGuests := 0
		notComingGuests := 0
		notAnsweredGuests := 0

		for i, guest := range msg.guests {
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
						comingGuests++
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
				strconv.FormatBool(guest.PlusOne()),
				strconv.FormatBool(guest.Coming()),
				"https://wedding.1995parham.ir/guests/" + guest.ID,
				strconv.FormatBool(guest.Answer == nil && !guest.IsFamily),
			}
		}

		m.isLoading = false
		m.table.SetRows(rows)

		m.text.SetValue(
			fmt.Sprintf(`
- Not Answered: %d / %d
- Coming (includes their plus one adult and their children): %d
- Not Coming: %d`,
				notAnsweredGuests,
				len(msg.guests),
				comingGuests,
				notComingGuests,
			),
		)
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m guestsModel) View() string {
	if m.isLoading {
		return fmt.Sprintf("\n\n   %s Loading from database...\n\n", m.spinner.View())
	}

	return m.table.View() + "\n\n\n" + m.text.View() + "\n"
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, repository guestrepo.Repository) {
	// nolint: mnd
	columns := []table.Column{
		{Title: "First Name", Width: 15},
		{Title: "Last Name", Width: 15},
		{Title: "Partner's First Name", Width: 15},
		{Title: "Partner's Last Name", Width: 15},
		{Title: "Family", Width: 10},
		{Title: "Children", Width: 10},
		{Title: "ID", Width: 10},
		{Title: "PlusOne", Width: 10},
		{Title: "Coming", Width: 10},
		{Title: "Link", Width: 50},
		{Title: "Waiting for an answer", Width: 15},
	}

	// nolint: mnd
	dm := guestsModel{
		repository: repository,
		isLoading:  true,
		spinner:    spinner.New(),
		table: table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
			table.WithHeight(50),
		),
		text: textarea.New(),
	}

	// nolint: mnd
	dm.text.SetWidth(200)
	dm.text.ShowLineNumbers = false

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
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "list",
		Description: "List all guests",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.NopLogger,
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(db.Provide),
				fx.Provide(
					fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
				),
				fx.Provide(generator.Provide),
				fx.Provide(service.ProvideGuestSvc),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
