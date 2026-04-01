package insert

import (
	"context"
	"errors"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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

var (
	errPartnerLastNameRequired = errors.New("partner last name is required when first name is provided")
	errMustBeNumber            = errors.New("must be a number")
	errCannotBeNegative        = errors.New("cannot be negative")
)

func weddingTheme() *huh.Theme {
	t := huh.ThemeCharm()

	orange := lipgloss.Color("#FF8C00")

	t.Focused.Title = t.Focused.Title.Foreground(orange)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Background(orange)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(orange)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(orange)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(orange)
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(orange)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(orange)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(orange)
	t.Group.Title = t.Group.Title.Foreground(orange)

	softOrange := lipgloss.Color("#FFB347")

	t.Blurred.Title = t.Blurred.Title.Foreground(softOrange)
	t.Blurred.TextInput.Placeholder = t.Blurred.TextInput.Placeholder.Foreground(lipgloss.Color("#997040"))
	t.Blurred.SelectedOption = t.Blurred.SelectedOption.Foreground(softOrange)
	t.Blurred.FocusedButton = t.Blurred.FocusedButton.Background(softOrange)
	t.Blurred.NextIndicator = t.Blurred.NextIndicator.Foreground(softOrange)

	return t
}

type guestForm struct {
	firstName        string
	lastName         string
	partnerFirstName string
	partnerLastName  string
	isFamily         bool
	children         string
}

func (gf *guestForm) build() *huh.Form { //nolint: funlen
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("First Name").
				Placeholder("Ali").
				Value(&gf.firstName).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Last Name").
				Placeholder("Irani").
				Value(&gf.lastName).
				Validate(huh.ValidateNotEmpty()),
		).Title("Guest Information"),

		huh.NewGroup(
			huh.NewInput().
				Title("Partner's First Name").
				Description("Leave empty if no partner").
				Placeholder("Maryam").
				Value(&gf.partnerFirstName),
			huh.NewInput().
				Title("Partner's Last Name").
				Placeholder("Akhyani").
				Value(&gf.partnerLastName).
				Validate(func(s string) error {
					if gf.partnerFirstName != "" && s == "" {
						return errPartnerLastNameRequired
					}

					return nil
				}),
		).Title("Partner Information"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Is this a family invitation?").
				Affirmative("Yes").
				Negative("No").
				Value(&gf.isFamily),
			huh.NewInput().
				Title("Number of Children").
				Placeholder("0").
				Value(&gf.children).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}

					n, err := strconv.Atoi(s)
					if err != nil {
						return errMustBeNumber
					}

					if n < 0 {
						return errCannotBeNegative
					}

					return nil
				}),
		).Title("Family Details"),
	).WithTheme(weddingTheme()).
		WithProgramOptions(tea.WithAltScreen())
}

func (gf *guestForm) numChildren() int {
	if gf.children == "" {
		return 0
	}

	n, _ := strconv.Atoi(gf.children)

	return n
}

func main(lc fx.Lifecycle, shutdowner fx.Shutdowner, svc service.GuestSvc) {
	gf := &guestForm{
		firstName:        "",
		lastName:         "",
		partnerFirstName: "",
		partnerLastName:  "",
		isFamily:         false,
		children:         "",
	}
	form := gf.build()

	lc.Append(
		fx.StartHook(func(ctx context.Context) error {
			go func() {
				defer func() { _ = shutdowner.Shutdown() }()

				if err := form.Run(); err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						return
					}

					pterm.Fatal.Printfln("form failed: %s", err)

					return
				}

				guest, err := svc.New(
					ctx,
					gf.firstName,
					gf.lastName,
					gf.partnerFirstName,
					gf.partnerLastName,
					gf.isFamily,
					gf.numChildren(),
				)
				if err != nil {
					pterm.Error.Printfln("failed to create guest: %s", err)

					return
				}

				pterm.Success.Printfln("Guest created: %s %s (ID: %s)", guest.FirstName, guest.LastName, guest.ID)

				if guest.SpouseFirstName != nil {
					pterm.Info.Printfln("Partner: %s %s", *guest.SpouseFirstName, *guest.SpouseLastName)
				}

				if gf.isFamily {
					pterm.Info.Printfln("Family invitation with %d children", gf.numChildren())
				}
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
