package settings

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	config "github.com/shinya-ac/server1Q1A/configs"
	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			c.Error(err)
			switch e := err.(type) {
			case *errDomain.Error:
				if errors.Is(e, errDomain.NotFoundErr) {
					return ReturnNotFound(c, e)
				}
				return ReturnStatusBadRequest(c, e)
			default:
				return ReturnStatusInternalServerError(c, e)
			}
		}
		return nil
	}
}

func ApiKeyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKeys := []string{config.Config.APIKey1, config.Config.APIKey2, config.Config.APIKey3}
		apiKey := c.Request().Header.Get("Todo-API-Key")

		valid := false
		for _, key := range apiKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "APIkeyが有効ではありません。"})
		}

		return next(c)
	}
}
