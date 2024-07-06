package health_handler

import (
	"github.com/labstack/echo/v4"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
)

func HealthCheck(ctx echo.Context) error {
	res := HealthResponse{
		Status: "ok",
	}
	return settings.ReturnStatusOK(ctx, res)
}
