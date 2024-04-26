package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"go.uber.org/zap"
)

type Guest struct {
	repository guestrepo.Repository
	logger     *zap.Logger
}

func (h Guest) Page(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	guest, err := h.repository.Get(ctx, id)
	if err != nil {
		h.logger.Error("failed to fetch a guest from repository", zap.Error(err), zap.String("id", id))

		if errors.Is(err, guestrepo.ErrGuestNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, guest)
}
