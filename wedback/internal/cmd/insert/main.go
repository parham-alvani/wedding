package insert

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type guestModel struct {
	service service.GuestSvc
	inputs  []textinput.Model
	index   int
}

func (m guestModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m guestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		// nolint: exhaustive
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.index == 0 {
				m.index++
			} else {
				if _, err := m.service.New(context.Background(), m.inputs[0].Value(), m.inputs[1].Value()); err != nil {
					pterm.Error.Printfln("failed to create the guest %s", err)
				}

				return m, tea.Quit
			}
		}
	}

	m.inputs[m.index], cmd = m.inputs[m.index].Update(msg)

	return m, cmd
}

func (m guestModel) View() string {
	return fmt.Sprintf(
		"What is your guest name?\n\n%s\n\nWhat is his/her partner name (Leave empty if there isn't any)?\n\n%s\n\n%s",
		m.inputs[0].View(),
		m.inputs[1].View(),
		"(esc to quit)",
	) + "\n"
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, svc service.GuestSvc) {
	iName := textinput.New()
	iName.Placeholder = "Ali Irani"
	iName.Focus()
	iName.CharLimit = 128
	iName.Width = 20

	iPartner := textinput.New()
	iPartner.Placeholder = "Maryam Akhyani"
	iPartner.Focus()
	iPartner.CharLimit = 128
	iPartner.Width = 20

	m := guestModel{
		service: svc,
		inputs:  []textinput.Model{iName, iPartner},
		index:   0,
	}

	p := tea.NewProgram(m)

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

// Register insert command.
func Register() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "insert",
		Description: "Insert a new guest",
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
