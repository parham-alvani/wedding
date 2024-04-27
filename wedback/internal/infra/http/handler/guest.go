package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/request"
	"go.uber.org/zap"
)

type Guest struct {
	Service service.GuestSvc
	Logger  *zap.Logger
}

func (h Guest) Page(c echo.Context) error {
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

	return c.JSON(http.StatusOK, guest)
}

func (h Guest) Answer(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var req request.Answer

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := h.Service.Answer(
		ctx,
		id,
		req.PlusOne,
		req.Coming,
	); err != nil {
		h.Logger.Error("failed to add an answer to a guest from repository", zap.Error(err), zap.String("id", id))

		if errors.Is(err, guestrepo.ErrGuestNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (h Guest) Register(g *echo.Group) {
	g.POST("/guest/:id/answer", h.Answer)
	g.GET("/guest/:id", h.Page)
}
