package settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoEngine() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	return e
}

func ReturnStatusOK[T any](ctx echo.Context, body T) error {
	return ctx.JSON(http.StatusOK, body)
}

func ReturnStatusCreated[T any](ctx echo.Context, body T) error {
	return ctx.JSON(http.StatusCreated, body)
}

func ReturnStatusNoContent(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func ReturnStatusBadRequest(ctx echo.Context, err error) error {
	return returnAbortWith(ctx, http.StatusBadRequest, err)
}

func ReturnBadRequest(ctx echo.Context, err error) error {
	return ReturnStatusBadRequest(ctx, err)
}

func ReturnStatusUnauthorized(ctx echo.Context, err error) error {
	return returnAbortWith(ctx, http.StatusUnauthorized, err)
}

func ReturnUnauthorized(ctx echo.Context, err error) error {
	return ReturnStatusUnauthorized(ctx, err)
}

func ReturnStatusForbidden(ctx echo.Context, err error) error {
	return returnAbortWith(ctx, http.StatusForbidden, err)
}

func ReturnForbidden(ctx echo.Context, err error) error {
	return ReturnStatusForbidden(ctx, err)
}

func ReturnStatusNotFound(ctx echo.Context, err error) error {
	return returnAbortWith(ctx, http.StatusNotFound, err)
}

func ReturnNotFound(ctx echo.Context, err error) error {
	return ReturnStatusNotFound(ctx, err)
}

func ReturnStatusInternalServerError(ctx echo.Context, err error) error {
	return returnAbortWith(ctx, http.StatusInternalServerError, err)
}

func ReturnError(ctx echo.Context, err error) error {
	ctx.Error(err)
	return nil
}

func returnAbortWith(ctx echo.Context, code int, err error) error {
	var msg string
	if err != nil {
		msg = err.Error()
	}

	return ctx.JSON(code, echo.Map{
		"code": code,
		"msg":  msg,
	})
}
