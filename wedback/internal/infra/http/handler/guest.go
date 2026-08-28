package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/request"
	"go.uber.org/zap"
)

type Guest struct {
	Service service.GuestSvc
	Logger  *zap.Logger
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
}
