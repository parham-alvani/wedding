package insert

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}

type guestModel struct {
	repository guestrepo.Repository
	input      textinput.Model
}

func (m guestModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m guestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if err := m.repository.Create(context.Background(), model.Guest{
				Name: m.input.Value(),
				ID:   RandStringRunes(10),
			}); err != nil {
				pterm.Error.Printfln("failed to create the guest %s", err)
			}

			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m guestModel) View() string {
	return fmt.Sprintf(
		"What is your guest name?\n\n%s\n\n%s",
		m.input.View(),
		"(esc to quit)",
	) + "\n"
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, repository guestrepo.Repository) {
	ti := textinput.New()
	ti.Placeholder = "Name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := guestModel{
		repository: repository,
		input:      ti,
	}

	p := tea.NewProgram(m)

	lc.Append(
		fx.StartHook(func(_ context.Context) error {
			if _, err := p.Run(); err != nil {
				return err
			}

			return shutdowner.Shutdown()
		}),
	)
}

// Register insert command.
func Register() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "insert",
		Description: "Insert a new guest",
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
