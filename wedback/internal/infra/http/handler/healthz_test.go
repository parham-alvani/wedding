package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/handler"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

type HealthzSuite struct {
	suite.Suite

	engine *echo.Echo
}

func (suite *HealthzSuite) SetupSuite() {
	suite.engine = echo.New()

	fxtest.New(suite.T(),
		fx.Invoke(func() {
			url := handler.Healthz{
				Logger: zap.NewNop(),
			}
			url.Register(suite.engine.Group(""))
		}),
	).RequireStart().RequireStop()
}

func (suite *HealthzSuite) TestHandler() {
	require := suite.Require()

	w := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.engine.ServeHTTP(w, req)
	require.Equal(http.StatusNoContent, w.Code)
}

func TestHealthzSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HealthzSuite))
}
