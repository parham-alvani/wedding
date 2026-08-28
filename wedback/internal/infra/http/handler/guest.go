package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/request"
	"github.com/parham-alvani/wedding/wedback/internal/infra/ratelimit"
	"go.uber.org/zap"
)

// verifyLimit bounds guesses at a guest's name. Generous enough that a guest
// fumbling their own spelling is never locked out, tight enough that working
// through common surnames is not worth trying.
const (
	verifyLimit  = 10
	verifyWindow = 10 * time.Minute
)

type Guest struct {
	Service service.GuestSvc
	Logger  *zap.Logger

	// attempts bounds name guesses per invitation.
	attempts *ratelimit.Limiter
}

// NewGuest builds the guest handler with its attempt limiter.
func NewGuest(svc service.GuestSvc, logger *zap.Logger) Guest {
	return Guest{
		Service:  svc,
		Logger:   logger,
		attempts: ratelimit.New(verifyLimit, verifyWindow),
	}
}

// Verify checks the name a visitor typed against the invitation. The reply
// says only yes or no: it never echoes the guest's real name back.
func (h Guest) Verify(c *echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var req request.Verify

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if h.attempts != nil && !h.attempts.Allow(id) {
		return echo.NewHTTPError(http.StatusTooManyRequests, "too many attempts, please try again later")
	}

	ok, err := h.Service.VerifyName(ctx, id, req.Name)
	if err != nil {
		if errors.Is(err, guestrepo.ErrGuestNotFound) {
			return echo.ErrNotFound
		}

		h.Logger.Error("failed to verify a guest name", zap.Error(err), zap.String("id", id))

		return echo.ErrInternalServerError
	}

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "that name does not match this invitation")
	}

	if h.attempts != nil {
		h.attempts.Reset(id)
	}

	return c.JSON(http.StatusOK, map[string]bool{"ok": true}) // nolint: wrapcheck
}

func (h Guest) Page(c *echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	guest, err := h.Service.Get(ctx, id)
	if err != nil {
		h.Logger.Error("failed to fetch a guest from repository", zap.Error(err), zap.String("id", id))

		if errors.Is(err, guestrepo.ErrGuestNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, guest) // nolint: wrapcheck
}

func (h Guest) Answer(c *echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var req request.Answer

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := h.Service.Answer(ctx, id, service.Reply{
		Coming:  req.Coming,
		PlusOne: req.PlusOne,
		Dietary: req.Dietary,
		Song:    req.Song,
	}); err != nil {
		h.Logger.Error("failed to add an answer to a guest from repository", zap.Error(err), zap.String("id", id))

		if errors.Is(err, guestrepo.ErrGuestNotFound) {
			return echo.ErrNotFound
		}

		// The guest did nothing wrong, they are simply too late; say so
		// rather than returning an opaque 500.
		if errors.Is(err, service.ErrRSVPClosed) {
			return echo.NewHTTPError(http.StatusForbidden, service.ErrRSVPClosed.Error())
		}

		if errors.Is(err, service.ErrNoteTooLong) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil) // nolint: wrapcheck
}

// Deadline reports whether the RSVP is closed and when it closes, so the
// frontend can render the invitation accordingly.
func (h Guest) Deadline(c *echo.Context) error {
	resp := map[string]any{"closed": h.Service.RSVPClosed()}

	if deadline := h.Service.Deadline(); !deadline.IsZero() {
		resp["deadline"] = deadline.Format(time.RFC3339)
	}

	return c.JSON(http.StatusOK, resp) // nolint: wrapcheck
}

func (h Guest) Register(g *echo.Group) {
	g.POST("/guest/:id/answer", h.Answer)
	g.GET("/guest/:id", h.Page)
	g.GET("/rsvp", h.Deadline)
	g.POST("/guest/:id/verify", h.Verify)
}
