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
	prompts []string
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
		case tea.KeyTab:
			m.index = (m.index + 1) % len(m.inputs)
		case tea.KeyShiftTab:
			m.index = (m.index - 1) % len(m.inputs)
		case tea.KeyEnter:
			if m.index < len(m.inputs)-1 {
				m.index++
			} else {
				if _, err := m.service.New(
					context.Background(),
					m.inputs[0].Value(),
					m.inputs[1].Value(),
					m.inputs[2].Value(),
					m.inputs[3].Value(),
				); err != nil {
					pterm.Error.Printfln("failed to create the guest %s", err)
				}

				return m, tea.Quit
			}
		}
	}

	for i := range len(m.inputs) {
		if i == m.index {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	m.inputs[m.index], cmd = m.inputs[m.index].Update(msg)

	return m, cmd
}

func (m guestModel) View() string {
	view := ""

	for i := range len(m.inputs) {
		view += fmt.Sprintf(
			"%s\n\n%s\n\n",
			m.prompts[i],
			m.inputs[i].View(),
		) + "\n"
	}

	view += "(esc to quit)"

	return view + "\n"
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, svc service.GuestSvc) {
	fName := textinput.New()
	fName.Placeholder = "Ali"
	fName.Focus()
	fName.CharLimit = 20
	fName.Width = 20

	lName := textinput.New()
	lName.Placeholder = "Irani"
	lName.CharLimit = 20
	lName.Width = 20

	fPartner := textinput.New()
	fPartner.Placeholder = "Maryam"
	fPartner.CharLimit = 20
	fPartner.Width = 20

	lPartner := textinput.New()
	lPartner.Placeholder = "Akhyani"
	lPartner.CharLimit = 20
	lPartner.Width = 20

	m := guestModel{
		service: svc,
		inputs:  []textinput.Model{fName, lName, fPartner, lPartner},
		prompts: []string{"First Name", "Last Name", "Partner's First Name", "Partner's Last Name"},
		index:   0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

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
		Aliases:     []string{"new"},
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
