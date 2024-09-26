package auth

import (
	"context"
	"errors"

	"github.com/form3tech-oss/jwt-go"
	"github.com/shinya-ac/server1Q1A/middlewares/auth0"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	token := auth0.GetJWT(ctx)
	if token == nil {
		logging.Logger.Error("JWT token not found in context")
		return "", errors.New("token not found")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logging.Logger.Error("Failed to cast token claims to jwt.MapClaims")
		return "", errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		logging.Logger.Error("sub claim not found in token claims")
		return "", errors.New("sub not found")
	}

	return sub, nil
}
