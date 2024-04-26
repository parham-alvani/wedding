package list

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type guestsModel struct {
	repository guestrepo.Repository

	isLoading bool

	table   table.Model
	spinner spinner.Model
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

func (m guestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg == nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case guestsListMsg:
		rows := make([]table.Row, len(msg.guests))
		for i, guest := range msg.guests {
			rows[i] = table.Row{
				guest.Name,
				guest.ID,
				fmt.Sprintf("%t", guest.PlusOne()),
				fmt.Sprintf("%t", guest.Coming()),
			}
		}

		m.isLoading = false
		m.table.SetRows(rows)
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m guestsModel) View() string {
	if m.isLoading {
		return fmt.Sprintf("\n\n   %s Loading from database...\n\n", m.spinner.View())
	}

	return baseStyle.Render(m.table.View()) + "\n"
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, repository guestrepo.Repository) {
	// nolint: gomnd
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "ID", Width: 10},
		{Title: "PlusOne", Width: 10},
		{Title: "Coming", Width: 10},
	}

	// nolint: gomnd
	dm := guestsModel{
		repository: repository,
		isLoading:  true,
		spinner:    spinner.New(),
		table: table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
			table.WithHeight(7),
		),
	}

	p := tea.NewProgram(dm)

	lc.Append(
		fx.StartHook(func(_ context.Context) error {
			if _, err := p.Run(); err != nil {
				return err
			}

			return shutdowner.Shutdown()
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
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(db.Provide),
				fx.Provide(
					fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
				),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
